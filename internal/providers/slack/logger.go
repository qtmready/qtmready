package slack

import (
	"log/slog"
)

type (
	logger struct {
		*slog.Logger
	}
)

func (l *logger) Output(calldepth int, s string) error {
	l.Logger.Info(s)

	return nil
}
