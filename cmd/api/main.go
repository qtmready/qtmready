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
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"

	"go.breu.io/quantm/internal/auth"
	"go.breu.io/quantm/internal/core"
	"go.breu.io/quantm/internal/providers/github"
	"go.breu.io/quantm/internal/shared"
	"go.breu.io/quantm/internal/shared/logger"
)

const (
	HTTPPort       = "8000"
	PrometheusPort = "9090"
)

func main() {
	exitcode := 0
	ctx := context.Background()
	quit := make(chan os.Signal, 1) // create a channel to listen to quit signals.
	errs := make(chan error, 1)     // create a channel to listen to errors.

	// init service
	shared.Service().SetName("api")
	shared.Logger().Info(
		"starting ...",
		slog.Any("service", shared.Service().GetName()),
		slog.String("version", shared.Service().GetVersion()),
	)

	// init open telemetry
	otelshutdown, err := _otel(
		ctx, shared.Service().GetName(),
		shared.Service().GetVersion(),
	)
	if err != nil {
		slog.Error("failed to setup opentelemetry, exiting ...", slog.Any("error", err.Error()))
		errs <- err

		return
	}

	slog.Info("setting up webserver")

	web := echo.New()
	web.HideBanner = true
	web.HTTPErrorHandler = shared.EchoAPIErrorHandler
	web.Validator = &shared.EchoValidator{Validator: shared.Validator()}

	web.Use(middleware.CORS())
	web.Use(otelecho.Middleware(shared.Service().GetName()))
	web.Use(logger.NewRequestLoggerMiddleware())
	web.Use(echoprometheus.NewMiddleware(shared.Service().GetName()))
	web.Use(middleware.Recover())

	web.GET("/healthz", healthz)

	auth.RegisterHandlers(web, auth.NewServerHandler(auth.Middleware))
	core.RegisterHandlers(web, core.NewServerHandler(auth.Middleware))
	github.RegisterHandlers(web, github.NewServerHandler(auth.Middleware))

	slog.Info("setting up metrics")

	metrics := echo.New()
	metrics.HideBanner = true

	// configure metrics routes
	metrics.GET("/metrics", echoprometheus.NewHandler())

	go _run(_serve(web, HTTPPort), errs)
	go _run(_serve(metrics, PrometheusPort), errs)
	slog.Info("registering quit signals")
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // setting up the signals to listen to.

	shared.Service().Banner()

	select {
	case err := <-errs:
		slog.Error("encountered error, shutting down gracefully ...", slog.Any("error", err.Error()))
		_graceful(ctx, []shutdownfn{web.Shutdown, metrics.Shutdown, otelshutdown}, []chan any{}, exitcode)

	case <-quit:
		slog.Info("received quit signal, shutting down gracefully ...")
		_graceful(ctx, []shutdownfn{web.Shutdown, metrics.Shutdown, otelshutdown}, []chan any{}, exitcode)
	}
}

// func main() {
// 	shared.Service().SetName("api")
// 	// graceful shutdown.
// 	// LINK: https://stackoverflow.com/a/46255965/228697.
// 	exitcode := 0
// 	defer func() {
// 		shared.Logger().Info("exiting ...")
// 		os.Exit(exitcode)
// 	}() // all connections are closed, exit with the right code.
// 	defer func() { _ = shared.Logger().Sync() }() // flush log buffer.
// 	defer db.DB().Session.Close()
// 	defer shared.Temporal().Client().Close()

// 	shared.Logger().Info("starting ...")
// 	// web server based on echo
// 	e := echo.New()
// 	e.HideBanner = true

// 	// configure middleware
// 	e.Use(middleware.CORS())
// 	e.Use(middleware.Logger())
// 	e.Use(middleware.Recover())
// 	e.Use(echoprometheus.NewMiddleware(shared.Service().GetName()))

// 	// e.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
// 	// 	body := string(reqBody[:])
// 	// 	fmt.Printf("bdy: %v", body)
// 	// }))

// 	// override the defaults
// 	e.Validator = &shared.EchoValidator{Validator: shared.Validator()}
// 	e.HTTPErrorHandler = shared.EchoAPIErrorHandler

// 	// register handlers
// 	auth.RegisterHandlers(e, auth.NewServerHandler(auth.Middleware))
// 	core.RegisterHandlers(e, core.NewServerHandler(auth.Middleware))
// 	github.RegisterHandlers(e, github.NewServerHandler(auth.Middleware))

// 	e.GET("/healthz", healthz)
// 	e.GET("/metrics", echoprometheus.NewHandler())

// 	go func() {
// 		if err := e.Start(":8000"); err != nil && err != http.ErrServerClosed {
// 			exitcode = 1
// 			return
// 		}
// 	}()

// 	shared.Service().Banner()

// 	quit := make(chan os.Signal, 1)                      // create a channel to listen to quit signals.
// 	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // setting up the signals to listen to.
// 	<-quit                                               // wait for quit signal.

// 	if err := e.Shutdown(context.Background()); err != nil {
// 		exitcode = 1
// 		return
// 	}

// 	exitcode = 1
// }
