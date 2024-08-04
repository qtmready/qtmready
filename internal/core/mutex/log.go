package mutex

import (
	"fmt"
	"strings"

	"go.temporal.io/sdk/log"
	"go.temporal.io/sdk/workflow"
)

type (
	// MutexLoggerKind defines the type of mutex logger.
	MutexLoggerKind string

	// MutexLogger provides logging functionality for mutex operations.
	MutexLogger struct {
		kind     MutexLoggerKind
		mutex_id string
		logger   log.Logger
	}

	// LogWriter defines a function type for writing log messages.
	LogWriter func(msg string, keyvals ...any)
)

const (
	// MutexHandlerKind represents the handler context for logging.
	MutexHandlerKind MutexLoggerKind = "mutex_hndl"
	// MutexControllerKind represents the controller context for logging.
	MutexControllerKind MutexLoggerKind = "mutex_ctrl"
)

// Logging methods

// Info logs an info-level message.
func (m *MutexLogger) info(handler_id, action, msg string, keyvals ...any) {
	m.write(m.logger.Info, handler_id, action, msg, keyvals...)
}

// Warn logs a warning-level message.
func (m *MutexLogger) warn(handler_id, action, msg string, keyvals ...any) {
	m.write(m.logger.Warn, handler_id, action, msg, keyvals...)
}

// Error logs an error-level message.
func (m *MutexLogger) error(handler_id, action, msg string, keyvals ...any) {
	m.write(m.logger.Error, handler_id, action, msg, keyvals...)
}

// Debug logs a debug-level message.
func (m *MutexLogger) debug(handler_id, action, msg string, keyvals ...any) {
	m.write(m.logger.Debug, handler_id, action, msg, keyvals...)
}

// Helper methods

// prefix creates a formatted prefix for log messages.
func (m *MutexLogger) prefix(handler_id, action string) string {
	return fmt.Sprintf("%s/%s/%s/%s: ", m.kind, m.strip(handler_id), m.strip(m.mutex_id), action)
}

// write handles the actual writing of log messages.
func (m *MutexLogger) write(writer LogWriter, handler_id, action, msg string, keyvals ...any) {
	keyvals = append(keyvals, "mutex_id", m.strip(m.mutex_id))
	keyvals = append(keyvals, "handler_id", m.strip(handler_id))
	keyvals = append(keyvals, "action", action)

	writer(m.prefix(handler_id, action)+msg, keyvals...)
}

// strip removes the first three parts of the ID if they exist.
func (m *MutexLogger) strip(id string) string {
	parts := strings.Split(id, ".")
	if len(parts) > 3 {
		return strings.Join(parts[3:], ".")
	}

	return id
}

// New methods

// NewMutexLogger creates a new MutexLogger instance.
func NewMutexLogger(ctx workflow.Context, kind MutexLoggerKind, mutex_id string) *MutexLogger {
	logger := workflow.GetLogger(ctx)

	return &MutexLogger{
		kind:     kind,
		mutex_id: mutex_id,
		logger:   logger,
	}
}

// NewMutexHandlerLogger creates a new MutexLogger instance for the handler context.
func NewMutexHandlerLogger(ctx workflow.Context, mutex_id string) *MutexLogger {
	return NewMutexLogger(ctx, MutexHandlerKind, mutex_id)
}

// NewMutexControllerLogger creates a new MutexLogger instance for the controller context.
func NewMutexControllerLogger(ctx workflow.Context, mutex_id string) *MutexLogger {
	return NewMutexLogger(ctx, MutexControllerKind, mutex_id)
}
