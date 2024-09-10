package testdata

import (
	"./shared"
)

func singleton_error() {
	shared.Logger().Error("error")
}
