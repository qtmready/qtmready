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
	Logger interface {
		Debug(string, ...any)
		DebugContext(context.Context, string, ...any)
		Error(string, ...any)
		ErrorContext(context.Context, string, ...any)
		Info(string, ...any)
		InfoContext(context.Context, string, ...any)
		Printf(string, ...any)
		Sync() error
		Trace(string, ...any)
		TraceContext(context.Context, string, ...any)
		Warn(string, ...any)
		WarnContext(context.Context, string, ...any)
	}

	// zapadapter is a wrapper around zap.Logger. Makes it compatible with the logger.Logger interface.
	zapadapter struct {
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

func NewZapAdapter(debug bool, skip int) Logger {
	logger := NewZapLogger(debug, skip)

	return &zapadapter{
		logger: logger, // skip the caller of this function
		core:   logger.Core(),
	}
}

func (a *zapadapter) Trace(msg string, fields ...any) {
	// TODO: Implement OpenTelemetry compatible Trace
	a.logger.Debug(msg, a.fields(fields)...)
}

func (a *zapadapter) Debug(msg string, fields ...any) {
	if !a.core.Enabled(zapcore.DebugLevel) {
		return
	}

	a.logger.Debug(msg, a.fields(fields)...)
}

func (a *zapadapter) Info(msg string, fields ...any) {
	if !a.core.Enabled(zapcore.InfoLevel) {
		return
	}

	a.logger.Info(msg, a.fields(fields)...)
}

func (a *zapadapter) Warn(msg string, fields ...any) {
	if !a.core.Enabled(zapcore.WarnLevel) {
		return
	}

	a.logger.Warn(msg, a.fields(fields)...)
}

func (a *zapadapter) Error(msg string, fields ...any) {
	if !a.core.Enabled(zapcore.ErrorLevel) {
		return
	}

	a.logger.Error(msg, a.fields(fields)...)
}

func (a *zapadapter) Printf(msg string, fields ...any) {
	a.logger.Sugar().Infof(msg, fields...)
}

func (a *zapadapter) Sync() error {
	return a.logger.Sync()
}

func (a *zapadapter) TraceContext(_ context.Context, msg string, fields ...any) {
	a.Trace(msg, fields...)
}

func (a *zapadapter) DebugContext(_ context.Context, msg string, fields ...any) {
	a.Debug(msg, fields...)
}

func (a *zapadapter) InfoContext(_ context.Context, msg string, fields ...any) {
	a.Info(msg, fields...)
}

func (a *zapadapter) WarnContext(_ context.Context, msg string, fields ...any) {
	a.Warn(msg, fields...)
}

func (a *zapadapter) ErrorContext(_ context.Context, msg string, fields ...any) {
	a.Error(msg, fields...)
}

func (a *zapadapter) fields(kv []any) []zap.Field {
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
