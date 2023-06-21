package testdata

import (
	"./shared"
)

var (
	logger = shared.Logger()
)

func context_info() {
	logger.Info("info")
}
