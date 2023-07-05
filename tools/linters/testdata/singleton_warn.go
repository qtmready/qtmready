package testdata

import (
	"./shared"
)

func singleton_warn() {
	shared.Logger().Warn("warn")
}
