package common

import (
	"go.uber.org/zap"
)

var (
	EventStream eventstreamconf
	Logger      *zap.Logger
	Service     serviceconf
	Temporal    temporalconf
)
