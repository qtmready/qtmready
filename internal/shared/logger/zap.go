// Copyright © 2022, Breu Inc. <info@breu.io>. All rights reserved. 
//
// This software is made available by Breu, Inc., under the terms of the 
// Breu Community License Agreement ("BCL Agreement"), version 1.0, found at  
// https://www.breu.io/license/community. By installating, downloading, 
// accessing, using or distrubting any of the software, you agree to the  
// terms of the license agreement. 
//
// The above copyright notice and the subsequent license agreement shall be 
// included in all copies or substantial portions of the software. 
//
// Breu, Inc. HEREBY DISCLAIMS ANY AND ALL WARRANTIES AND CONDITIONS, EXPRESS, 
// IMPLIED, STATUTORY, OR OTHERWISE, AND SPECIFICALLY DISCLAIMS ANY WARRANTY OF 
// MERCHANTABILITY OR FITNESS FOR A PARTICULAR PURPOSE, WITH RESPECT TO THE 
// SOFTWARE. 
//
// Breu, Inc. SHALL NOT BE LIABLE FOR ANY DAMAGES OF ANY KIND, INCLUDING BUT NOT 
// LIMITED TO, LOST PROFITS OR ANY CONSEQUENTIAL, SPECIAL, INCIDENTAL, INDIRECT, 
// OR DIRECT DAMAGES, HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, ARISING 
// OUT OF THIS AGREEMENT. THE FOREGOING SHALL APPLY TO THE EXTENT PERMITTED BY  
// APPLICABLE LAW. 

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
