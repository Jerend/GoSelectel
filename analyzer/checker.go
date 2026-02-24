package analyzer

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"
	"unicode"

	"golang.org/x/tools/go/analysis"
)

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			if callExpr, ok := n.(*ast.CallExpr); ok {
				checkLogCall(pass, callExpr)
			}
			return true
		})
	}
	return nil, nil
}

func checkLogCall(pass *analysis.Pass, call *ast.CallExpr) {
	if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
		if pkg, ok := sel.X.(*ast.Ident); ok && pkg.Name == "slog" {
			checkSlogCall(pass, call, sel.Sel.Name)
			return
		}
	}

	var ident *ast.Ident
	var findBaseIdent func(expr ast.Expr) bool
	findBaseIdent = func(expr ast.Expr) bool {
		switch e := expr.(type) {
		case *ast.Ident:
			ident = e
			return true
		case *ast.SelectorExpr:
			return findBaseIdent(e.X)
		case *ast.CallExpr:
			return findBaseIdent(e.Fun)
		default:
			return false
		}
	}

	if !findBaseIdent(call.Fun) {
		return
	}

	var methodName string
	if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
		methodName = sel.Sel.Name
	} else {
		return
	}

	if ident != nil {
		obj := pass.TypesInfo.ObjectOf(ident)
		if obj == nil {
			return
		}
		objType := obj.Type().String()

		switch {
		case strings.Contains(objType, "log/slog.Logger"):
			checkSlogCall(pass, call, methodName)
		case strings.Contains(objType, "zap.Logger"):
			checkZapCall(pass, call, methodName)
		case strings.Contains(objType, "zap.SugaredLogger"):
			checkZapSugaredCall(pass, call, methodName)
		}
	}
}

func checkSlogCall(pass *analysis.Pass, call *ast.CallExpr, method string) {
	switch method {
	case "Debug", "Info", "Warn", "Error":
		if len(call.Args) > 0 {
			checkLogMessage(pass, call.Args[0])
			for i := 0; i < len(call.Args); i++ {
				checkSensitiveData(pass, call.Args[i])
			}
		}
	}
}

func checkZapCall(pass *analysis.Pass, call *ast.CallExpr, method string) {
	switch method {
	case "Debug", "Info", "Warn", "Error", "DPanic", "Panic", "Fatal":
		if len(call.Args) == 0 {
			return
		}
		checkLogMessage(pass, call.Args[0])
		checkSensitiveData(pass, call.Args[0])
		for i := 1; i < len(call.Args); i++ {
			if field, ok := call.Args[i].(*ast.CallExpr); ok {
				if sel, ok := field.Fun.(*ast.SelectorExpr); ok {
					if pkg, ok := sel.X.(*ast.Ident); ok && pkg.Name == "zap" {
						if len(field.Args) >= 2 {
							checkSensitiveData(pass, field.Args[1])
						}
						if sel.Sel.Name == "Error" && len(field.Args) >= 1 {
							checkSensitiveData(pass, field.Args[0])
						}
					}
				}
			}
		}
	}
}

func checkZapSugaredCall(pass *analysis.Pass, call *ast.CallExpr, method string) {
	if len(call.Args) == 0 {
		return
	}

	hasSuffix := func(s string) bool {
		return strings.HasSuffix(method, s)
	}

	switch {
	case hasSuffix("f"): // методы "Debugf", "Infof", "Warnf", "Errorf", "DPanicf", "Panicf", "Fatalf"
		if len(call.Args) >= 1 {
			checkLogMessage(pass, call.Args[0])
			for i := 1; i < len(call.Args); i++ {
				checkSensitiveData(pass, call.Args[i])
			}
		}
	case hasSuffix("w"): //методы ключ-значнеие "Debugw", "Infow", "Warnw", "Errorw", "DPanicw", "Panicw", "Fatalw"
		if len(call.Args) >= 1 {
			checkLogMessage(pass, call.Args[0])
			for i := 1; i < len(call.Args); i++ {
				if i%2 == 0 {
					checkSensitiveData(pass, call.Args[i])
				}
			}
		}
	default: //методы "Debug", "Info", "Warn", "Error", "DPanic", "Panic", "Fatal"
		switch method {
		case "Debug", "Info", "Warn", "Error", "DPanic", "Panic", "Fatal":
			if len(call.Args) >= 1 {
				checkLogMessage(pass, call.Args[0])
				checkSensitiveData(pass, call.Args[0])
				for i := 1; i < len(call.Args); i++ {
					checkSensitiveData(pass, call.Args[i])
				}
			}
		}
	}
}

func checkLogMessage(pass *analysis.Pass, messArg ast.Expr) {
	var msg string
	var pos, end token.Pos
	var originalText string

	switch arg := messArg.(type) {
	case *ast.BasicLit:
		if arg.Kind != token.STRING {
			return
		}
		originalText = arg.Value
		msg = strings.Trim(originalText, "\"`")
		pos = arg.Pos()
		end = arg.End()
	case *ast.BinaryExpr:
		checkLogMessage(pass, arg.X)
		checkLogMessage(pass, arg.Y)
		return
	default:
		return
	}

	if len(msg) == 0 {
		return
	}

	firstChar := []rune(msg)[0]
	if unicode.IsUpper(firstChar) {
		fixed := string(unicode.ToLower(firstChar)) + msg[1:]
		fixedMsg := strings.Replace(originalText, msg, fixed, 1)

		pass.Report(analysis.Diagnostic{
			Pos:     pos,
			Message: fmt.Sprintf("The message in the log starts with a capital letter, got %q", msg),
			SuggestedFixes: []analysis.SuggestedFix{{
				Message: "Convert first letter to lowercase",
				TextEdits: []analysis.TextEdit{{
					Pos:     pos,
					End:     end,
					NewText: []byte(fixedMsg),
				}},
			}},
		})
	}

	var firstNonEnglishIdx = -1
	runes := []rune(msg)

	for i, r := range runes {
		if unicode.IsLetter(r) && !((r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z')) {
			firstNonEnglishIdx = i
			break
		}
	}

	if firstNonEnglishIdx != -1 {
		var clean strings.Builder
		for _, r := range runes {
			if !unicode.IsLetter(r) || ((r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z')) {
				clean.WriteRune(r)
			}
		}

		fixedMsg := strings.Replace(originalText, msg, clean.String(), 1)

		pass.Report(analysis.Diagnostic{
			Pos:     pos,
			Message: "log message contains non-English characters",
			SuggestedFixes: []analysis.SuggestedFix{{
				Message: "Remove non-English characters",
				TextEdits: []analysis.TextEdit{{
					Pos:     pos,
					End:     end,
					NewText: []byte(fixedMsg),
				}},
			}},
		})
	}

	allowedSpecial := map[rune]bool{
		',': true, ':': true, ';': true, '-': true,
		'/': true, '(': true, ')': true, '[': true, ']': true,
		'{': true, '}': true, '=': true, '+': true, ' ': true,
	}

	var invalidPositions []int
	runes = []rune(msg)
	for i, r := range runes {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && !allowedSpecial[r] {
			invalidPositions = append(invalidPositions, i)
		}
	}

	if len(invalidPositions) > 0 {
		var clean strings.Builder
		for i, r := range runes {
			shouldKeep := true
			for _, idx := range invalidPositions {
				if i == idx {
					shouldKeep = false
					break
				}
			}
			if shouldKeep {
				clean.WriteRune(r)
			}
		}

		fixedMsg := strings.Replace(originalText, msg, clean.String(), 1)

		pass.Report(analysis.Diagnostic{
			Pos:     pos,
			Message: "log message contains invalid characters",
			SuggestedFixes: []analysis.SuggestedFix{{
				Message: "Remove invalid characters",
				TextEdits: []analysis.TextEdit{{
					Pos:     pos,
					End:     end,
					NewText: []byte(fixedMsg),
				}},
			}},
		})
	}
}

func checkSensitiveData(pass *analysis.Pass, expr ast.Expr) {
	if binaryExpr, ok := expr.(*ast.BinaryExpr); ok && binaryExpr.Op == token.ADD {
		checkSensitiveData(pass, binaryExpr.X)
		checkSensitiveData(pass, binaryExpr.Y)
		return
	}

	if callExpr, ok := expr.(*ast.CallExpr); ok {
		if sel, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
			if sel.Sel.Name == "Sprintf" || strings.HasSuffix(sel.Sel.Name, "printf") {
				for i := 1; i < len(callExpr.Args); i++ {
					checkSensitiveData(pass, callExpr.Args[i])
				}
			}
		}
		return
	}

	if ident, ok := expr.(*ast.Ident); ok {
		name := strings.ToLower(ident.Name)
		sensitive := []string{"password", "pwd", "pass", "apikey", "token", "auth_token", "secret", "secret_key", "private_key", "access_key", "login"}

		for _, sens := range sensitive {
			if strings.Contains(name, sens) {
				pass.Report(analysis.Diagnostic{
					Pos:     ident.Pos(),
					Message: fmt.Sprintf("potential sensitive data %q is being logged directly", ident.Name),
					SuggestedFixes: []analysis.SuggestedFix{{
						Message: "Redact sensitive data",
						TextEdits: []analysis.TextEdit{{
							Pos:     ident.Pos(),
							End:     ident.End(),
							NewText: []byte(`"[change]"`),
						}},
					}},
				})
				return
			}
		}
	}
}
