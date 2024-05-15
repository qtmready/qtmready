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
	"io"
	"log/slog"

	"cloud.google.com/go/compute/metadata"
	"go.opentelemetry.io/otel/trace"
)

type (
	GoogleCloudHandler struct {
		handler slog.Handler
	}
)

func NewGoogleCloudHandler(writer io.Writer, options *slog.HandlerOptions) slog.Handler {
	options.ReplaceAttr = replaceattr

	handler := slog.NewJSONHandler(writer, options)

	return &GoogleCloudHandler{handler: handler}
}

func (h *GoogleCloudHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *GoogleCloudHandler) Handle(ctx context.Context, rec slog.Record) error {
	return h.handler.Handle(ctx, h.enrich(ctx, rec))
}

func (h *GoogleCloudHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &GoogleCloudHandler{handler: h.handler.WithAttrs(attrs)}
}

func (h *GoogleCloudHandler) WithGroup(name string) slog.Handler {
	return &GoogleCloudHandler{handler: h.handler.WithGroup(name)}
}

// enrich adds the trace ID to the record so it is correlated with the Cloud Run request log
//
// # LINKS
//   - https://cloud.google.com/trace/docs/trace-log-integration
//   - https://cloud.google.com/logging/docs/view/correlate-logs#view-correlated-log-entries
func (h *GoogleCloudHandler) enrich(ctx context.Context, record slog.Record) slog.Record {
	rec := record.Clone()

	span := trace.SpanFromContext(ctx)
	if span != nil && span.SpanContext().IsValid() {
		if metadata.OnGCE() {
			project, _ := metadata.ProjectID()
			rec.Add("logging.googleapis.com/trace", fmt.Sprintf(
				"projects/%s/traces/%s",
				project,
				span.SpanContext().TraceID().String(),
			))
		} else {
			rec.Add("logging.googleapis.com/trace", span.SpanContext().TraceID().String())
		}

		rec.Add("logging.googleapis.com/spanId", span.SpanContext().SpanID().String())
		rec.Add("logging.googleapis.com/trace_sampled", span.SpanContext().IsSampled())
	}

	return rec
}

func replaceattr(groups []string, attr slog.Attr) slog.Attr {
	switch attr.Key {
	case slog.MessageKey:
		attr.Key = "message"

	case slog.SourceKey:
		attr.Key = "logging.googleapis.com/sourceLocation"

	case slog.TimeKey:
		attr.Key = "timestamp"

	case slog.LevelKey:
		attr.Key = "severity"
		level := attr.Value.Any().(slog.Level)

		if level == slog.LevelWarn {
			attr.Value = slog.StringValue("WARNING")
		}
	}

	return attr
}
