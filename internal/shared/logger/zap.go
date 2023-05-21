// Copyright © 2023, Breu, Inc. <info@breu.io>. All rights reserved.
//
// This software is made available by Breu, Inc., under the terms of the BREU COMMUNITY LICENSE AGREEMENT, Version 1.0,
// found at https://www.breu.io/license/community. BY INSTALLING, DOWNLOADING, ACCESSING, USING OR DISTRIBUTING ANY OF
// THE SOFTWARE, YOU AGREE TO THE TERMS OF THE LICENSE AGREEMENT.
//
// The above copyright notice and the subsequent license agreement shall be included in all copies or substantial
// portions of the software.
//
// Breu, Inc. HEREBY DISCLAIMS ANY AND ALL WARRANTIES AND CONDITIONS, EXPRESS, IMPLIED, STATUTORY, OR OTHERWISE, AND
// SPECIFICALLY DISCLAIMS ANY WARRANTY OF MERCHANTABILITY OR FITNESS FOR A PARTICULAR PURPOSE, WITH RESPECT TO THE
// SOFTWARE.
//
// Breu, Inc. SHALL NOT BE LIABLE FOR ANY DAMAGES OF ANY KIND, INCLUDING BUT NOT LIMITED TO, LOST PROFITS OR ANY
// CONSEQUENTIAL, SPECIAL, INCIDENTAL, INDIRECT, OR DIRECT DAMAGES, HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY,
// ARISING OUT OF THIS AGREEMENT. THE FOREGOING SHALL APPLY TO THE EXTENT PERMITTED BY APPLICABLE LAW.

package logger

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type (
	// ZapAdapter is a wrapper around zap.Logger. Makes it compatible with the logger.Logger interface.
	ZapAdapter struct {
		logger *zap.Logger
		core   zapcore.Core
	}
)

func NewZapLogger(debug bool, skip int) *zap.Logger {
	var config zap.Config
	if debug {
		config = zap.NewDevelopmentConfig()
	} else {
		config = zap.NewProductionConfig()
	}

	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	config.EncoderConfig.CallerKey = "func"

	zl, _ := config.Build(zap.AddCallerSkip(skip))

	return zl
}

func NewZapAdapter(logger *zap.Logger, skip int) *ZapAdapter {
	return &ZapAdapter{
		logger: logger.WithOptions(zap.AddCallerSkip(skip)), // skip the caller of this function
		core:   logger.Core(),
	}
}

func (adapter *ZapAdapter) Trace(msg string, fields ...any) {
	// TODO: Implement OpenTelemetry compatible Trace
	adapter.logger.Debug(msg, adapter.fields(fields)...)
}

func (adapter *ZapAdapter) Debug(msg string, fields ...any) {
	if !adapter.core.Enabled(zapcore.DebugLevel) {
		return
	}

	adapter.logger.Debug(msg, adapter.fields(fields)...)
}

func (adapter *ZapAdapter) Info(msg string, fields ...any) {
	if !adapter.core.Enabled(zapcore.InfoLevel) {
		return
	}

	adapter.logger.Info(msg, adapter.fields(fields)...)
}

func (adapter *ZapAdapter) Warn(msg string, fields ...any) {
	if !adapter.core.Enabled(zapcore.WarnLevel) {
		return
	}

	adapter.logger.Warn(msg, adapter.fields(fields)...)
}

func (adapter *ZapAdapter) Error(msg string, fields ...any) {
	if !adapter.core.Enabled(zapcore.ErrorLevel) {
		return
	}

	adapter.logger.Error(msg, adapter.fields(fields)...)
}

func (adapter *ZapAdapter) Printf(msg string, fields ...any) {
	adapter.logger.Sugar().Infof(msg, fields...)
}

func (adapter *ZapAdapter) Sync() error {
	return adapter.logger.Sync()
}

func (adapter *ZapAdapter) TraceContext(_ context.Context, msg string, fields ...any) {
	adapter.Trace(msg, fields...)
}

func (adapter *ZapAdapter) DebugContext(_ context.Context, msg string, fields ...any) {
	adapter.Debug(msg, fields...)
}

func (adapter *ZapAdapter) InfoContext(_ context.Context, msg string, fields ...any) {
	adapter.Info(msg, fields...)
}

func (adapter *ZapAdapter) WarnContext(_ context.Context, msg string, fields ...any) {
	adapter.Warn(msg, fields...)
}

func (adapter *ZapAdapter) ErrorContext(_ context.Context, msg string, fields ...any) {
	adapter.Error(msg, fields...)
}

func (adapter *ZapAdapter) fields(kv []any) []zap.Field {
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
