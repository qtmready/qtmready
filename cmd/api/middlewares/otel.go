package middlewares

import (
	"net/http"
	"sync"

	"github.com/felixge/httpsnoop"
	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/contrib"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

var tracerName = "go.breu.io/ctrlplane/cmd/api/middlewares"

// Chi middleware for OpenTelemetry
func OtelMiddleware(service string, options ...Option) func(next http.Handler) http.Handler {
	cfg := &conf{}

	for _, option := range options {
		option.apply(cfg)
	}

	if cfg.TracerProvider == nil {
		cfg.TracerProvider = otel.GetTracerProvider()
	}

	if cfg.Propagators == nil {
		cfg.Propagators = otel.GetTextMapPropagator()
	}

	tracer := cfg.TracerProvider.Tracer(tracerName, trace.WithInstrumentationVersion(contrib.SemVersion()))

	return func(handler http.Handler) http.Handler {
		return TraceOption{
			service:             service,
			tracer:              tracer,
			propagators:         cfg.Propagators,
			handler:             handler,
			routes:              cfg.Routes,
			reqMethodInSpanName: cfg.RequestMethodInSpanName,
		}
	}
}

type conf struct {
	TracerProvider          trace.TracerProvider
	Propagators             propagation.TextMapPropagator
	Routes                  chi.Routes
	RequestMethodInSpanName bool
}

type Option interface {
	apply(*conf)
}

type optionFn func(*conf)

func (fn optionFn) apply(cfg *conf) {
	fn(cfg)
}

func IncludeOtelPropagators(propagators propagation.TextMapPropagator) Option {
	return optionFn(func(cfg *conf) {
		cfg.Propagators = propagators
	})
}

func IncludeOtelTraceProviders(provider trace.TracerProvider) Option {
	return optionFn(func(cfg *conf) {
		cfg.TracerProvider = provider
	})
}

func IncludeChiRoutes(routes chi.Routes) Option {
	return optionFn(func(cfg *conf) {
		cfg.Routes = routes
	})
}

func IncludeRequestMethodInSpanName(isActive bool) Option {
	return optionFn(func(cfg *conf) {
		cfg.RequestMethodInSpanName = isActive
	})
}

type ResponseWriterRecorder struct {
	status     int
	isRecorded bool
	writer     http.ResponseWriter
}

var rwrPool = &sync.Pool{
	New: func() interface{} {
		return &ResponseWriterRecorder{}
	},
}

func getRWR(response http.ResponseWriter) *ResponseWriterRecorder {
	rwr := rwrPool.Get().(*ResponseWriterRecorder)
	rwr.isRecorded = false
	rwr.status = 0
	rwr.writer = httpsnoop.Wrap(response, httpsnoop.Hooks{
		WriteHeader: func(next httpsnoop.WriteHeaderFunc) httpsnoop.WriteHeaderFunc {
			return func(status int) {
				if !rwr.isRecorded {
					rwr.isRecorded = true
					rwr.status = status
				}
				next(status)
			}
		},
		Write: func(next httpsnoop.WriteFunc) httpsnoop.WriteFunc {
			return func(p []byte) (int, error) {
				if !rwr.isRecorded {
					rwr.isRecorded = true
					rwr.status = http.StatusOK
				}
				return next(p)
			}
		},
	})
	return rwr
}

func putRWR(rwr *ResponseWriterRecorder) {
	rwr.writer = nil
	rwrPool.Put(rwr)
}

type TraceOption struct {
	service             string
	tracer              trace.Tracer
	propagators         propagation.TextMapPropagator
	handler             http.Handler
	routes              chi.Routes
	reqMethodInSpanName bool
}

func (opt TraceOption) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	ctx := opt.propagators.Extract(request.Context(), propagation.HeaderCarrier(request.Header))
	name := ""
	pattern := ""

	// Find the route that matches the request and set the name and pattern accordingly.
	if opt.routes != nil {
		rctx := chi.NewRouteContext()
		if opt.routes.Match(rctx, request.Method, request.URL.Path) {
			pattern = rctx.RoutePattern()
			name = addMethodToSpan(opt.reqMethodInSpanName, request.Method, pattern)
		}
	}

	// Starting a new trace
	ctx, span := opt.tracer.Start(
		ctx,
		name,
		trace.WithAttributes(semconv.NetAttributesFromHTTPRequest("tcp", request)...),
		trace.WithAttributes(semconv.EndUserAttributesFromHTTPRequest(request)...),
		trace.WithAttributes(semconv.HTTPServerAttributesFromHTTPRequest(opt.service, pattern, request)...),
	)

	defer span.End()

	// Record the span in the context
	rwr := getRWR(response)
	defer putRWR(rwr)

	request = request.WithContext(ctx)

	// Handle the next request
	opt.handler.ServeHTTP(rwr.writer, request)

	// Setting span attributes if required
	if len(pattern) == 0 {
		pattern = chi.RouteContext(request.Context()).RoutePattern()
		span.SetAttributes(semconv.HTTPRouteKey.String(pattern))
		name = addMethodToSpan(opt.reqMethodInSpanName, request.Method, pattern)
		span.SetName(name)
	}

	// Setting HTTP status code on the span
	span.SetAttributes(semconv.HTTPStatusCodeKey.Int(rwr.status))
	status, msg := semconv.SpanStatusFromHTTPStatusCode(rwr.status)
	span.SetStatus(status, msg)
}

// If the request method is not included in the span name, add it to the span name.
func addMethodToSpan(shouldAdd bool, method string, span string) string {
	if shouldAdd {
		span = method + span
	}
	return span
}
