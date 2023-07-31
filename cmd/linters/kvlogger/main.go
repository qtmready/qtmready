package main

import (
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/singlechecker"

	"go.breu.io/quantm/tools/linters"
)

var (
	AnalyzerPlugin plugin
)

type (
	plugin struct{}
)

func (p *plugin) GetAnalyzers() []*analysis.Analyzer {
	return []*analysis.Analyzer{
		linters.KVLoggerAnalyzer,
	}
}

func main() {
	singlechecker.Main(linters.KVLoggerAnalyzer)
}

func New(conf any) ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{
		linters.KVLoggerAnalyzer,
	}, nil
}
