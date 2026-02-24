package unittests

import (
	"GolangTestSelectel/analyzer"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestEnglish_only(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), analyzer.Analyzer, "english_only")
}
