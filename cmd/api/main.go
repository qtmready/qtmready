// Copyright © 2022, Breu, Inc. <info@breu.io>. All rights reserved.
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
	"net/http"
	"os"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	swagger "github.com/swaggo/echo-swagger"

	"go.breu.io/ctrlplane/cmd/api/docs"
	"go.breu.io/ctrlplane/internal/auth"
	"go.breu.io/ctrlplane/internal/core"
	"go.breu.io/ctrlplane/internal/db"
	"go.breu.io/ctrlplane/internal/providers"
	"go.breu.io/ctrlplane/internal/providers/github"
	"go.breu.io/ctrlplane/internal/shared"
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
	waigroup := sync.WaitGroup{}
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
	shared.Logger.Info("initializing ...")
	waigroup.Add(3)

	go func() {
		defer waigroup.Done()
		db.DB.InitSession()
	}()

	go func() {
		defer waigroup.Done()
		shared.EventStream.InitConnection()
	}()

	go func() {
		defer waigroup.Done()
		shared.Temporal.InitClient()
	}()

	waigroup.Wait()
	// Initializing singleton objects ... Done

	shared.Logger.Info("initialized", "version", shared.Service.Version())
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
