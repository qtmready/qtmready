package status

import (
	"go.breu.io/quantm/internal/erratic"
)

const (
	ConnectionError = 200001
)

func NewConnectionError(args ...string) *erratic.QuantmError {
	return erratic.New(ConnectionError, "Connection Error.", args...)
}
