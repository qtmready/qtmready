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

package main

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"cloud.google.com/go/compute/metadata"
	cloudtrace "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	gcppropagator "github.com/GoogleCloudPlatform/opentelemetry-operations-go/propagator"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

type (
	shutdownfn func(ctx context.Context) error

	devnull struct{}
)

func (d *devnull) ExportSpans(ctx context.Context, spans []trace.ReadOnlySpan) error {
	return nil
}

func (d *devnull) Shutdown(ctx context.Context) error {
	return nil
}

func (d *devnull) MarshalLog() any {
	return nil
}

// _otel sets up opentelemetry.
func _otel(ctx context.Context, name, version string) (shutdown shutdownfn, err error) {
	slog.Info("otel: init ...")

	shutdownfns := make([]shutdownfn, 0)
	shutdown = func(ctx context.Context) error {
		var errs error
		for _, fn := range shutdownfns {
			errs = errors.Join(errs, fn(ctx))
		}

		shutdownfns = nil

		return errs
	}

	handlerr := func(err error) {
		err = errors.Join(err, shutdown(ctx)) // nolint: ineffassign, staticcheck
	}

	/**
	 * Setup OpenTelemetry Exporter
	 */

	exporter, err := _exporter()
	if err != nil {
		handlerr(err)
		return
	}

	/**
	 * Setup OpenTelemetry Resource
	 */

	res, err := _resource(ctx, name, version)
	if err != nil {
		handlerr(err)
		return
	}

	/**
	 * Setup OpenTelemetry Propagator
	 */

	propagator := _propagator()
	otel.SetTextMapPropagator(propagator)

	/**
	 * Setup OpenTelemetry Trace Provider
	 */

	tracer := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithBatcher(exporter, trace.WithBatchTimeout(time.Second)),
		trace.WithResource(res),
	)

	otel.SetTracerProvider(tracer)

	shutdownfns = append(shutdownfns, tracer.Shutdown)

	slog.Info("otel: initialized")

	return
}

// _exporter returns a trace exporter based on the environment. When running on GCE, the cloudtrace exporter is used.
func _exporter() (trace.SpanExporter, error) {
	if metadata.OnGCE() {
		project, _ := metadata.ProjectID()
		return cloudtrace.New(cloudtrace.WithProjectID(project))
	}

	// return stdouttrace.New()
	return &devnull{}, nil
}

// _resource returns a resource with the service name and version.
func _resource(_ context.Context, name, version string) (*resource.Resource, error) {
	return resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL, semconv.ServiceName(name), semconv.ServiceVersion(version)),
	)
}

// _propagator returns a propagator based on the environment. The gcp.CloudTraceOneWayPropagator is only used when
// running on GCE. It allows for the open telemetry traceheader to take precedence over the x-cloud-trace-context.
func _propagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		gcppropagator.CloudTraceOneWayPropagator{},
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}
