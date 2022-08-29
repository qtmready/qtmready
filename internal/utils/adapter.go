package utils

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Adapter struct {
	logger *zap.Logger
	core   zapcore.Core
}

func NewZapAdapter(logger *zap.Logger) *Adapter {
	return &Adapter{
		logger: logger.WithOptions(zap.AddCallerSkip(1)), // skip the caller of this function
		// logger: logger,
		core: logger.Core(),
	}
}

func (adapter *Adapter) Trace(msg string, fields ...interface{}) {
	// TODO: Implement OpenTelemetry compatible Trace
	adapter.logger.Debug(msg, adapter.fields(fields)...)
}

func (adapter *Adapter) Debug(msg string, fields ...interface{}) {
	if !adapter.core.Enabled(zapcore.DebugLevel) {
		return
	}

	adapter.logger.Debug(msg, adapter.fields(fields)...)
}

func (adapter *Adapter) Info(msg string, fields ...interface{}) {
	if !adapter.core.Enabled(zapcore.InfoLevel) {
		return
	}

	adapter.logger.Info(msg, adapter.fields(fields)...)
}

func (adapter *Adapter) Warn(msg string, fields ...interface{}) {
	if !adapter.core.Enabled(zapcore.WarnLevel) {
		return
	}

	adapter.logger.Warn(msg, adapter.fields(fields)...)
}

func (adapter *Adapter) Error(msg string, fields ...interface{}) {
	if !adapter.core.Enabled(zapcore.ErrorLevel) {
		return
	}

	adapter.logger.Error(msg, adapter.fields(fields)...)
}

func (adapter *Adapter) TraceContext(_ context.Context, msg string, fields ...interface{}) {
	adapter.Trace(msg, fields...)
}

func (adapter *Adapter) DebugContext(_ context.Context, msg string, fields ...interface{}) {
	adapter.Debug(msg, fields...)
}

func (adapter *Adapter) InfoContext(_ context.Context, msg string, fields ...interface{}) {
	adapter.Info(msg, fields...)
}

func (adapter *Adapter) WarnContext(_ context.Context, msg string, fields ...interface{}) {
	adapter.Warn(msg, fields...)
}

func (adapter *Adapter) ErrorContext(_ context.Context, msg string, fields ...interface{}) {
	adapter.Error(msg, fields...)
}

func (adapter *Adapter) fields(keyvals []interface{}) []zap.Field {
	var fields []zap.Field
	if len(keyvals)%2 != 0 {
		for _, f := range keyvals {
			fields = append(fields, f.(zap.Field))
		}
		// return []zap.Field{zap.Error(fmt.Errorf("odd number of keyvals pairs: %v", keyvals))}
	}
	for i := 0; i < len(keyvals); i += 2 {
		key, ok := keyvals[i].(string)
		if !ok {
			key = fmt.Sprintf("%v", keyvals[i])
		}
		fields = append(fields, zap.Any(key, keyvals[i+1]))
	}

	return fields
}
