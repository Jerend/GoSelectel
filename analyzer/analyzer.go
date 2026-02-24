package analyzer

import "golang.org/x/tools/go/analysis"

var Analyzer = &analysis.Analyzer{
	Name: "mylinter",
	Doc:  "check log messages",
	Run:  run,
}
