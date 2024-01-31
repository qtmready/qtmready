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
	SimpleCloudRunHandler struct {
		handler slog.Handler
	}
)

func NewSimpleCloudRunHandler(writer io.Writer, options *slog.HandlerOptions) slog.Handler {
	options.ReplaceAttr = replaceattr

	handler := slog.NewJSONHandler(writer, options)

	return &SimpleCloudRunHandler{handler: handler}
}

func (h *SimpleCloudRunHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *SimpleCloudRunHandler) Handle(ctx context.Context, rec slog.Record) error {
	return h.handler.Handle(ctx, h.enrich(ctx, rec))
}

func (h *SimpleCloudRunHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &SimpleCloudRunHandler{handler: h.handler.WithAttrs(attrs)}
}

func (h *SimpleCloudRunHandler) WithGroup(name string) slog.Handler {
	return &SimpleCloudRunHandler{handler: h.handler.WithGroup(name)}
}

// enrich adds the trace ID to the record so it is correlated with the Cloud Run request log
//
// # LINKS
//   - https://cloud.google.com/trace/docs/trace-log-integration
//   - https://cloud.google.com/logging/docs/view/correlate-logs#view-correlated-log-entries
func (h *SimpleCloudRunHandler) enrich(ctx context.Context, record slog.Record) slog.Record {
	rec := record.Clone()

	span := trace.SpanFromContext(ctx)
	if span != nil {
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
