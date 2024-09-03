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
	"log/slog"
	"regexp"
	"runtime/debug"

	"cloud.google.com/go/compute/metadata"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type (
	CloudTraceContextKey string
)

const (
	CloudTraceContextHeader = "X-Cloud-Trace-Context"

	TraceContextKey   CloudTraceContextKey = "trace"
	SpanContextKey    CloudTraceContextKey = "span"
	SampledContextKey CloudTraceContextKey = "trace_sampled"
)

var (
	DefaultRequestLoggerConfig = middleware.RequestLoggerConfig{
		Skipper:          middleware.DefaultSkipper,
		LogValuesFunc:    EchoRequestLogger,
		HandleError:      true,
		LogContentLength: true,
		LogError:         true,
		LogHost:          true,
		LogLatency:       true,
		LogMethod:        true,
		LogProtocol:      true,
		LogReferer:       true,
		LogRemoteIP:      true,
		LogResponseSize:  true,
		LogStatus:        true,
		LogURI:           true,
		LogUserAgent:     true,
	}

	condition = regexp.MustCompile(
		// Matches on "TRACE_ID"
		`([a-f\d]+)?` +
			// Matches on "/SPAN_ID"
			`(?:/([a-f\d]+))?` +
			// Matches on ";0=TRACE_TRUE"
			`(?:;o=(\d))?`)
)

// ParseCloudTraceHeaderMiddleware is a middleware that parses the x-cloud-trace-context header and adds the trace, span and
// sampled values to the request context. We then pass the context to the logger so that the trace and span are
// included in the logs.
//
// This will also come in handy when we are instrumenting third party call with OpenTelemetry.
func ParseCloudTraceHeaderMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		if metadata.OnGCE() {
			enrich(ctx)
		}

		return next(ctx)
	}
}

// EchoRequestLogger is a middleware that logs the request.
// NOTE: This should come after the ParseCloudTraceMiddleware.
func EchoRequestLogger(ctx echo.Context, values middleware.RequestLoggerValues) error {
	level := slog.LevelInfo
	url := fmt.Sprintf("%s://%s%s", ctx.Scheme(), values.Host, values.URI)

	attrs := []slog.Attr{
		slog.Any("latency", values.Latency),
		slog.String("latency_human", fmt.Sprintf("%d %s", values.Latency.Milliseconds(), "ms")),
		slog.String("protocol", values.Protocol),
		slog.String("referer", values.Referer),
		slog.String("remote_ip", values.RemoteIP),
		slog.String("request_method", values.Method),
		slog.Int64("request_size", ctx.Request().ContentLength),
		slog.String("request_url", url),
		slog.Int64("response_size", values.ResponseSize),
		slog.Int("status", values.Status),
		slog.String("user_agent", values.UserAgent),
		slog.String("host", values.Host),
	}

	if values.Error != nil {
		level = slog.LevelError

		// known errors: logged as warning without stack trace
		// system or unhandled error: logged as error with stack trace
		if values.Status <= 499 && values.Status > 399 {
			level = slog.LevelWarn
		} else {
			var stack []byte
			stack = debug.Stack()
			attrs = append(attrs, slog.String("stack_trace", string(stack)))
		}

		slog.Default().LogAttrs(ctx.Request().Context(), level, values.Error.Error(), attrs...)

		return values.Error
	}

	if values.URI == "/healthx" {
		level = slog.LevelDebug
	}

	slog.Default().LogAttrs(ctx.Request().Context(), level, url, slog.Any("request", attrs))

	return values.Error
}

// NewRequestLoggerMiddleware returns a new instance of RequestLogger middleware with default settings.
func NewRequestLoggerMiddleware() echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(DefaultRequestLoggerConfig)
}

// enrich parses the x-cloud-trace-context header and adds the trace, span and sampled values to the request context.
func enrich(ctx echo.Context) echo.Context {
	header := ctx.Request().Header.Get(CloudTraceContextHeader)
	if header != "" && condition != nil {
		matches := condition.FindStringSubmatch(header)
		trace, span, sampled := matches[1], matches[2], matches[3] == "1"

		if span == "0" {
			span = ""
		}

		requestctx := ctx.Request().Context()
		requestctx = context.WithValue(requestctx, TraceContextKey, trace)
		requestctx = context.WithValue(requestctx, SpanContextKey, span)
		requestctx = context.WithValue(requestctx, SampledContextKey, sampled)

		ctx.SetRequest(ctx.Request().WithContext(requestctx))
	}

	return ctx
}
