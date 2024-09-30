// Crafted with ❤ at Breu, Inc. <info@breu.io>, Copyright © 2024.
//
// Functional Source License, Version 1.1, Apache 2.0 Future License
//
// We hereby irrevocably grant you an additional license to use the Software under the Apache License, Version 2.0 that
// is effective on the second anniversary of the date we make the Software available. On or after that date, you may use
// the Software under the Apache License, Version 2.0, in which case the following will apply:
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
// the License.
//
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
// an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.

package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"

	"cloud.google.com/go/compute/metadata"
	"go.opentelemetry.io/otel/trace"
)

const (
	prefix_trace   = "logging.googleapis.com/trace"
	prefix_span    = "logging.googleapis.com/spanId"
	prefix_sampled = "logging.googleapis.com/trace_sampled"
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
			project, _ := metadata.ProjectIDWithContext(ctx)
			rec.Add(prefix_trace, fmt.Sprintf("projects/%s/traces/%s", project, span.SpanContext().TraceID().String()))
		} else {
			rec.Add(prefix_trace, span.SpanContext().TraceID().String())
		}

		rec.Add(prefix_span, span.SpanContext().SpanID().String())
		rec.Add(prefix_sampled, span.SpanContext().IsSampled())
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
