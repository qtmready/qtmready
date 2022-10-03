// Copyright © 2022, Breu Inc. <info@breu.io>. All rights reserved.

package main

import (
	"net/http"
	"os"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

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
	// graceful shutdown. see https://stackoverflow.com/a/46255965/228697.
	exitcode := 0
	defer func() { os.Exit(exitcode) }()
	defer func() { _ = shared.Logger.Sync() }()       // flush log buffer
	defer func() { _ = shared.EventStream.Drain() }() // process events in the buffer before closing connection
	defer db.DB.Session.Close()
	defer shared.Temporal.Client.Close()

	e := echo.New()

	e.Validator = &EchoValidator{validator: shared.Validate}

	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// A Mix of public & authenticated routes
	e.GET("/healthcheck", healthcheck)
	auth.CreateRoutes(e.Group("/auth"))
	providers.CreateRoutes(e.Group("/providers"), auth.Middleware)

	// Private routes
	protected := e.Group("")
	protected.Use(auth.Middleware)
	core.CreateRoutes(protected)

	if err := e.Start(":8000"); err != nil {
		exitcode = 1
		return
	}
}

// healthcheck checks if the system is working properly.
//
// TODO: make sure that connection to all services it needs to connect to is working properly.
func healthcheck(ctx echo.Context) error { return ctx.String(http.StatusOK, "OK") }
