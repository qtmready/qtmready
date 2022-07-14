package main

import (
	"context"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"go.opentelemetry.io/otel"
	stdout "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	prop "go.opentelemetry.io/otel/propagation"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	"go.breu.io/ctrlplane/cmd/api/middlewares"
	cmn "go.breu.io/ctrlplane/internal/common"
	"go.breu.io/ctrlplane/internal/db"
	"go.breu.io/ctrlplane/internal/integrations"
	"go.breu.io/ctrlplane/internal/integrations/github"
)

var waiter sync.WaitGroup
var traceProvider *sdktrace.TracerProvider

func init() {
	cmn.Service.ReadEnv()
	cmn.Service.InitLogger()

	cmn.EventStream.ReadEnv()
	cmn.Temporal.ReadEnv()
	github.Github.ReadEnv()
	db.DB.ReadEnv()

	waiter.Add(4)

	go func() {
		defer waiter.Done()
		db.DB.InitSessionWithMigrations()
	}()

	go func() {
		defer waiter.Done()
		cmn.EventStream.InitConnection()
	}()

	go func() {
		defer waiter.Done()
		cmn.Temporal.InitClient()
	}()

	go func() {
		defer waiter.Done()
		traceProvider = initTraceProvider()
	}()

	waiter.Wait()

	cmn.Logger.Info("Initializing Service ... Done")
}

func main() {
	// handling closing of the server
	defer db.DB.Session.Close()
	defer cmn.Temporal.Client.Close()
	defer func() {
		if err := traceProvider.Shutdown(context.Background()); err != nil {
			cmn.Logger.Error(err.Error())
		}
	}()

	// Setting up OpenTelemetry Global Tracer
	otel.SetTracerProvider(traceProvider)
	otel.SetTextMapPropagator(
		prop.NewCompositeTextMapPropagator(prop.TraceContext{}, prop.Baggage{}),
	)

	router := chi.NewRouter()

	router.Use(chimw.RequestID)
	router.Use(chimw.RealIP)
	router.Use(chimw.Logger)
	router.Use(middlewares.OtelMiddleware(
		cmn.Service.Name,
		middlewares.IncludeChiRoutes(router),
	))
	router.Use(chimw.Recoverer)

	router.Mount("/integrations", integrations.Router())

	http.ListenAndServe(":8000", router)
}

func initTraceProvider() *sdktrace.TracerProvider {
	cmn.Logger.Info("Initializing OpenTelemetry Provider ... ")
	exporter, err := stdout.New(stdout.WithPrettyPrint())
	if err != nil {
		cmn.Logger.Fatal(err.Error())
	}

	resource, err := sdkresource.New(
		context.Background(),
		sdkresource.WithAttributes(semconv.ServiceNameKey.String(cmn.Service.Name)),
	)

	if err != nil {
		cmn.Logger.Fatal(err.Error())
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(resource),
		sdktrace.WithBatcher(exporter),
	)

	cmn.Logger.Info("Initializing OpenTelemetry Provider ... Done")
	return tp
}
