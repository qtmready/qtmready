// Copyright © 2022, Breu Inc. <info@breu.io>. All rights reserved.  

package main

import (
	"net/http"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"go.breu.io/ctrlplane/cmd/api/auth"
	"go.breu.io/ctrlplane/cmd/api/core"
	"go.breu.io/ctrlplane/internal/db"
	"go.breu.io/ctrlplane/internal/drivers"
	"go.breu.io/ctrlplane/internal/drivers/github"
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

	shared.Logger.Info("Initializing Service ... Done")
}

func main() {
	// handling closing of the server
	defer db.DB.Session.Close()
	defer shared.Temporal.Client.Close()
	defer shared.Logger.Sync()

	e := echo.New()
	jwtconf := middleware.JWTConfig{
		Claims:     &shared.JWTClaims{},
		SigningKey: []byte(shared.Service.Secret),
	}

	e.Validator = &EchoValidator{validator: shared.Validate}

	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Unauthenticated routes
	e.GET("/healthcheck", healthcheck)
	auth.CreateRoutes(e.Group("/auth"))
	drivers.CreateRoutes(e.Group("/drivers"), middleware.JWTWithConfig(jwtconf))

	// Protected routes
	protected := e.Group("")
	protected.Use(middleware.JWTWithConfig(jwtconf))
	core.CreateRoutes(protected)

	if err := e.Start(":8000"); err != nil {
		shared.Logger.Error("Error starting web server", "error", err)
	}
}

// TODO: ensure connectivity with external services
func healthcheck(ctx echo.Context) error { return ctx.String(http.StatusOK, "OK") }
