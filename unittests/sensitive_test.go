package unittests

import (
	"GolangTestSelectel/analyzer"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestSensitive(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), analyzer.Analyzer, "sensitive")
}
