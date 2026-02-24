package unittests

import (
	"GolangTestSelectel/analyzer"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestSpecial_chars(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), analyzer.Analyzer, "special_chars")
}
