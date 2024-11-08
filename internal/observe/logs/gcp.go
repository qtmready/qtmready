package logs

import (
	"context"
	"fmt"
	"io"
	"log/slog"

	"cloud.google.com/go/compute/metadata"
	"go.opentelemetry.io/otel/trace"
)

const (
	prefix_trace   = "logging.googleapis.com/trace"         // Key for trace ID in log record.
	prefix_span    = "logging.googleapis.com/spanId"        // Key for span ID in log record.
	prefix_sampled = "logging.googleapis.com/trace_sampled" // Key for trace sampling status in log record.
)

type (
	// GoogleCloudHandler is an slog handler which adapts slog json handler to conform to Google Cloud Logging.
	// It is meant to be used with stdout as io.Writer.
	GoogleCloudHandler struct {
		handler slog.Handler // Underlying slog.Handler to delegate to.
	}
)

// NewGoogleCloudHandler constructs a new GoogleCloudHandler.
//
// It uses the provided io.Writer to write logs, and applies the specified slog.HandlerOptions.
// It modifies the options to map standard keys to Google Cloud Logging keys, and uses a JSON formatter.
func NewGoogleCloudHandler(writer io.Writer, options *slog.HandlerOptions) slog.Handler {
	options.ReplaceAttr = replaceattr // Replace standard keys with Google Cloud Logging keys.

	handler := slog.NewJSONHandler(writer, options) // Create a JSON-based handler.

	return &GoogleCloudHandler{handler: handler} // Return a new GoogleCloudHandler.
}

// Enabled checks if the handler is enabled for the given context and level.
//
// It delegates to the underlying handler's Enabled method.
func (h *GoogleCloudHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

// Handle writes a log record to the underlying handler.
//
// It enriches the record with trace context from the provided context, and then delegates to the underlying handler's
// Handle method.
func (h *GoogleCloudHandler) Handle(ctx context.Context, rec slog.Record) error {
	return h.handler.Handle(ctx, h.enrich(ctx, rec))
}

// WithAttrs returns a new handler with the specified attributes appended to the existing ones.
//
// It delegates to the underlying handler's WithAttrs method.
func (h *GoogleCloudHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &GoogleCloudHandler{handler: h.handler.WithAttrs(attrs)}
}

// WithGroup returns a new handler with the specified group name.
//
// It delegates to the underlying handler's WithGroup method.
func (h *GoogleCloudHandler) WithGroup(name string) slog.Handler {
	return &GoogleCloudHandler{handler: h.handler.WithGroup(name)}
}

// enrich adds trace context to the record.
//
// If the context contains a valid trace.Span, it extracts the trace ID, span ID, and sampling status, and adds them to the
// record as attributes. It also uses the metadata package to determine if the application is running on Google Cloud.
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

// replaceattr replaces standard log keys with Google Cloud Logging keys.
//
// It maps the following keys:
//   - slog.MessageKey -> "message"
//   - slog.SourceKey -> "logging.googleapis.com/sourceLocation"
//   - slog.TimeKey -> "timestamp"
//   - slog.LevelKey -> "severity"
//
// It also converts the slog.Level to a string value for "severity," using "WARNING" for slog.LevelWarn.
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
