// Copyright © 2022, Breu Inc. <info@breu.io>. All rights reserved.

package logger

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ZapAdapter is a wrapper around zap.Logger. Makes it compatible with the logger.Logger interface.
type ZapAdapter struct {
	logger *zap.Logger
	core   zapcore.Core
}

func NewZapAdapter(logger *zap.Logger) *ZapAdapter {
	return &ZapAdapter{
		logger: logger.WithOptions(zap.AddCallerSkip(1)), // skip the caller of this function
		core:   logger.Core(),
	}
}

func (adapter *ZapAdapter) Trace(msg string, fields ...interface{}) {
	// TODO: Implement OpenTelemetry compatible Trace
	adapter.logger.Debug(msg, adapter.fields(fields)...)
}

func (adapter *ZapAdapter) Debug(msg string, fields ...interface{}) {
	if !adapter.core.Enabled(zapcore.DebugLevel) {
		return
	}

	adapter.logger.Debug(msg, adapter.fields(fields)...)
}

func (adapter *ZapAdapter) Info(msg string, fields ...interface{}) {
	if !adapter.core.Enabled(zapcore.InfoLevel) {
		return
	}

	adapter.logger.Info(msg, adapter.fields(fields)...)
}

func (adapter *ZapAdapter) Warn(msg string, fields ...interface{}) {
	if !adapter.core.Enabled(zapcore.WarnLevel) {
		return
	}

	adapter.logger.Warn(msg, adapter.fields(fields)...)
}

func (adapter *ZapAdapter) Error(msg string, fields ...interface{}) {
	if !adapter.core.Enabled(zapcore.ErrorLevel) {
		return
	}

	adapter.logger.Error(msg, adapter.fields(fields)...)
}

func (adapter *ZapAdapter) Sync() error {
	return adapter.logger.Sync()
}

func (adapter *ZapAdapter) TraceContext(_ context.Context, msg string, fields ...interface{}) {
	adapter.Trace(msg, fields...)
}

func (adapter *ZapAdapter) DebugContext(_ context.Context, msg string, fields ...interface{}) {
	adapter.Debug(msg, fields...)
}

func (adapter *ZapAdapter) InfoContext(_ context.Context, msg string, fields ...interface{}) {
	adapter.Info(msg, fields...)
}

func (adapter *ZapAdapter) WarnContext(_ context.Context, msg string, fields ...interface{}) {
	adapter.Warn(msg, fields...)
}

func (adapter *ZapAdapter) ErrorContext(_ context.Context, msg string, fields ...interface{}) {
	adapter.Error(msg, fields...)
}

func (adapter *ZapAdapter) fields(kv []interface{}) []zap.Field {
	var fields []zap.Field

	if len(kv)%2 != 0 {
		return []zap.Field{zap.Error(fmt.Errorf("odd number of kv pairs: %v", kv))}
	}

	for i := 0; i < len(kv); i += 2 {
		key, ok := kv[i].(string)
		if !ok {
			key = fmt.Sprintf("%v", kv[i])
		}

		fields = append(fields, zap.Any(key, kv[i+1]))
	}

	return fields
}
