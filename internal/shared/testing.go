// Copyright © 2022, Breu Inc. <info@breu.io>. All rights reserved.

package shared

import (
	"testing"
)

type (
	TestFn struct {
		Args interface{} // Can be nil
		Want interface{} // Can be nil
		Run  func(provide interface{}, want interface{}) func(*testing.T)
	}

	TestFnMap map[string]TestFn
)
