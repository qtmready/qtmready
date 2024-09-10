package testutils

import (
	"testing"
)

type (
	TestFn struct {
		Args any // Can be nil
		Want any // Can be nil
		Run  func(provide, want any) func(*testing.T)
	}

	TestFnMap map[string]TestFn
)
