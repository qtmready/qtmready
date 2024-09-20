// Crafted with ❤ at Breu, Inc. <info@breu.io>, Copyright © 2022, 2024.
//
// Functional Source License, Version 1.1, Apache 2.0 Future License
//
// We hereby irrevocably grant you an additional license to use the Software under the Apache License, Version 2.0 that
// is effective on the second anniversary of the date we make the Software available. On or after that date, you may use
// the Software under the Apache License, Version 2.0, in which case the following will apply:
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
// the License.
//
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
// an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.

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
	"go.breu.io/quantm/internal/core/defs"
	"go.breu.io/quantm/internal/core/kernel"
	coreweb "go.breu.io/quantm/internal/core/web"
	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/providers/github"
	"go.breu.io/quantm/internal/providers/slack"
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

	defer db.DB().Session.Close()

	// init service
	shared.Service().SetName("api")
	shared.Logger().Info(
		"starting ...",
		slog.Any("service", shared.Service().GetName()),
		slog.String("version", shared.Service().GetVersion()),
	)

	// init open telemetry
	otelshutdown, err := _otel(ctx, shared.Service().GetName(), shared.Service().GetVersion())
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

	web.GET("/healthx", healthz)

	auth.RegisterHandlers(web, auth.NewServerHandler(auth.Middleware))
	coreweb.RegisterHandlers(web, coreweb.NewServerHandler(auth.Middleware))
	github.RegisterHandlers(web, github.NewServerHandler(auth.Middleware))
	slack.RegisterHandlers(web, slack.NewServerHandler(auth.Middleware))

	kernel.Instance(
		kernel.WithRepoProvider(defs.RepoProviderGithub, &github.RepoIO{}),
		kernel.WithMessageProvider(defs.MessageProviderSlack, &slack.Activities{}),
	)

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
