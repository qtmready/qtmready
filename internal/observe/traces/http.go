package trace

import (
	"context"
	"net/http"
	"regexp"
)

type (
	// CloudTraceContextKey defines a type for the Cloud Trace Context key.
	CloudTraceContextKey string
)

const (
	// CloudTraceContextHeader is the header name for the Cloud Trace Context.
	CloudTraceContextHeader = "X-Cloud-Trace-Context"

	// TraceContextKey is the key for storing the trace ID in the request context.
	TraceContextKey CloudTraceContextKey = "trace"

	// SpanContextKey is the key for storing the span ID in the request context.
	SpanContextKey CloudTraceContextKey = "span"

	// SampledContextKey is the key for storing the trace sampling flag in the request context.
	SampledContextKey CloudTraceContextKey = "trace_sampled"
)

var (
	// match defines a regular expression for parsing the Cloud Trace Context header.
	match = regexp.MustCompile(
		// Matches on "TRACE_ID"
		`([a-f\d]+)?` +
			// Matches on "/SPAN_ID"
			`(?:/([a-f\d]+))?` +
			// Matches on ";0=TRACE_TRUE"
			`(?:;o=(\d))?`)
)

// CloudTraceReporter extracts trace, span, and sampling flags from a HTTP request
// header and stores them in the context.
//
// It parses the `X-Cloud-Trace-Context` header and stores the extracted
// values in the context using the `TraceContextKey`, `SpanContextKey`, and
// `SampledContextKey` keys, respectively.
func CloudTraceReporter(req *http.Request) *http.Request {
	header := req.Header.Get(CloudTraceContextHeader)
	if header != "" && match != nil {
		matches := match.FindStringSubmatch(header)
		trace, span, sampled := matches[1], matches[2], matches[3] == "1"

		if span == "0" {
			span = ""
		}

		ctx := req.Context()
		ctx = context.WithValue(ctx, TraceContextKey, trace)
		ctx = context.WithValue(ctx, SpanContextKey, span)
		ctx = context.WithValue(ctx, SampledContextKey, sampled)

		req = req.WithContext(ctx)
	}

	return req
}
