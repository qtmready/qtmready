package linters_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"go.breu.io/quantm/tools/linters"
)

func TestKVLogger(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), linters.KVLoggerAnalyzer)
}
