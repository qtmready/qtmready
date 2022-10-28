// Copyright © 2022, Breu Inc. <info@breu.io>. All rights reserved.

package main

import (
	"net/http"
	"os"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	swagger "github.com/swaggo/echo-swagger"

	"go.breu.io/ctrlplane/cmd/api/docs"
	"go.breu.io/ctrlplane/internal/api/auth"
	"go.breu.io/ctrlplane/internal/api/core"
	"go.breu.io/ctrlplane/internal/db"
	"go.breu.io/ctrlplane/internal/providers"
	"go.breu.io/ctrlplane/internal/providers/github"
	"go.breu.io/ctrlplane/internal/shared"
)

var (
	waiter sync.WaitGroup
)

type (
	EchoValidator struct {
		validator *validator.Validate
	}
)

func (ev *EchoValidator) Validate(i interface{}) error {
	if err := ev.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return nil
}

func init() {
	// Reading the configuration from the environment
	shared.Service.ReadEnv()
	shared.Service.InitLogger()
	shared.Service.InitValidator()
	shared.EventStream.ReadEnv()
	db.DB.ReadEnv()
	db.DB.RegisterValidations()
	shared.Temporal.ReadEnv()
	github.Github.ReadEnv()
	// Reading the configuration from the environment ... Done

	// Initializing reference to adapters
	waiter.Add(3)

	go func() {
		defer waiter.Done()
		db.DB.InitSessionWithMigrations()
	}()

	go func() {
		defer waiter.Done()
		shared.EventStream.InitConnection()
	}()

	go func() {
		defer waiter.Done()
		shared.Temporal.InitClient()
	}()

	waiter.Wait()
	// Initializing singleton objects ... Done

	shared.Logger.Info("Initializing Service ... Done", "version", shared.Service.Version())
}

func main() {
	// graceful shutdown.
	// LINK: https://stackoverflow.com/a/46255965/228697.
	exitcode := 0
	defer func() { os.Exit(exitcode) }()              // all connections are closed, exit with the right code
	defer func() { _ = shared.Logger.Sync() }()       // flush log buffer
	defer func() { _ = shared.EventStream.Drain() }() // process events in the buffer before closing connection
	defer db.DB.Session.Close()
	defer shared.Temporal.Client.Close()

	// docs
	docs.SwaggerInfo.Title = shared.Service.Name
	docs.SwaggerInfo.Version = shared.Service.Version()

	// web server based on echo
	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Validator = &EchoValidator{validator: shared.Validate}

	// Public endpoints
	e.GET("/docs/*", swagger.WrapHandler)
	e.GET("/healthcheck", healthcheck)
	// Auth endpoints
	auth.CreateRoutes(e.Group("/auth"))

	// endpoints for 3rd party providers
	providers.CreateRoutes(e.Group("/providers"), auth.Middleware)

	// core api endpoints
	protected := e.Group("", auth.Middleware)
	core.CreateRoutes(protected)

	if err := e.Start(":8000"); err != nil {
		exitcode = 1
		return
	}
}

type (
	HealthCheckResponse struct {
		Msg string `json:"msg"`
	}
)

// healthcheck is the health check endpoint.
//
// @Summary     Checks if connection to all external services are working fine.
// @Description Quick health check
// @Tags        healthcheck
// @Accept      json
// @Produce     json
// @Success     201 {object} HealthCheckResponse
// @Failure     500 {object} echo.HTTPError
// @Router      /healthcheck [get]
//
// TODO: make sure that connection to all services it needs to connect to is working properly.
func healthcheck(ctx echo.Context) error {
	return ctx.Bind(&HealthCheckResponse{Msg: "OK"})
}
