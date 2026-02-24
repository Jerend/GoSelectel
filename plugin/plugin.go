package main

import (
	"GolangTestSelectel/analyzer"

	"golang.org/x/tools/go/analysis"
)

const Name = "mylinter"

var AnalyzerPlugin analyzerPlugin

type analyzerPlugin struct{}

func (analyzerPlugin) GetAnalyzers() []*analysis.Analyzer {
	return []*analysis.Analyzer{
		analyzer.Analyzer,
	}
}
