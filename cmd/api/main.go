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
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sourcegraph/conc"

	"go.breu.io/ctrlplane/internal/auth"
	"go.breu.io/ctrlplane/internal/core"
	"go.breu.io/ctrlplane/internal/db"
	"go.breu.io/ctrlplane/internal/providers/github"
	"go.breu.io/ctrlplane/internal/shared"
)

func init() {
	waitgroup := conc.WaitGroup{}
	defer waitgroup.Wait()
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
	waitgroup.Go(db.DB.InitSession)
	waitgroup.Go(shared.EventStream.InitConnection)
	waitgroup.Go(shared.Temporal.InitClient)
	shared.Logger.Info("initialized", "version", shared.Service.Version())
}

func main() {
	// graceful shutdown.
	// LINK: https://stackoverflow.com/a/46255965/228697.
	exitcode := 0
	defer func() {
		shared.Logger.Debug("exiting ...")
		os.Exit(exitcode)
	}() // all connections are closed, exit with the right code.
	defer func() { _ = shared.Logger.Sync() }()       // flush log buffer.
	defer func() { _ = shared.EventStream.Drain() }() // process events in the buffer before closing connection.
	defer db.DB.Session.Close()
	defer shared.Temporal.Client.Close()

	shared.Logger.Debug("starting ...")
	// web server based on echo
	e := echo.New()

	// configure middleware
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
		body := string(reqBody[:])
		shared.Logger.Debug("body: %s", body)
	}))

	// adding prometheus metrics
	prom := prometheus.NewPrometheus(shared.Service.Name, nil)
	prom.Use(e)

	// override the defaults
	e.Validator = &shared.EchoValidator{Validator: shared.Validator}
	e.HTTPErrorHandler = shared.EchoAPIErrorHandler

	// register handlers
	auth.RegisterHandlers(e, auth.NewServerHandler(auth.Middleware))
	core.RegisterHandlers(e, core.NewServerHandler(auth.Middleware))
	github.RegisterHandlers(e, github.NewServerHandler(auth.Middleware))

	e.GET("/healthz", healthz)

	go func() {
		if err := e.Start(":8000"); err != nil && err != http.ErrServerClosed {
			exitcode = 1
			return
		}
	}()

	quit := make(chan os.Signal, 1)                                       // create a channel to listen for quit signals.
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL) // listen for quit signals.
	<-quit                                                                // wait for quit signal.

	if err := e.Shutdown(context.Background()); err != nil {
		exitcode = 1
		return
	}

	exitcode = 1
}

type (
	HealthzResponse struct {
		Status string `json:"status"`
	}
)

// healthz is the health check endpoint.
//
// TODO: make sure that connection to all services it needs to connect to are working properly.
func healthz(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, &HealthzResponse{Status: "OK"})
}
