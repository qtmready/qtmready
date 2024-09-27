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
	"time"

	"github.com/labstack/echo/v4"

	"go.breu.io/quantm/internal/core/defs"
	"go.breu.io/quantm/internal/core/kernel"
	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/providers/github"
	"go.breu.io/quantm/internal/providers/slack"
	"go.breu.io/quantm/internal/shared"
	"go.breu.io/quantm/internal/shared/graceful"
)

const (
	HTTPPort       = "8000"
	PrometheusPort = "9090"
)

func main() {
	ctx := context.Background()
	sigterm := make(chan os.Signal, 1) // create a channel to listen to interrupt signals.
	interrupt := make(chan any, 1)     // channel to signal the shutdown to goroutines.
	errs := make(chan error)           // create a channel to listen to errors.

	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, os.Interrupt)

	// init service
	shared.Service().SetName("api")
	shared.Logger().Info(
		"starting ...",
		slog.Any("service", shared.Service().GetName()),
		slog.String("version", shared.Service().GetVersion()),
	)

	kernel.Instance(
		kernel.WithRepoProvider(defs.RepoProviderGithub, &github.RepoIO{}),
		kernel.WithMessageProvider(defs.MessageProviderSlack, &slack.Activities{}),
	)

	otelshutdown, err := observe(ctx, shared.Service().GetName(), shared.Service().GetVersion())
	if err != nil {
		slog.Error("failed to setup opentelemetry, exiting ...", slog.Any("error", err.Error()))
		errs <- err

		return
	}

	slog.Info("setting up webserver")

	web := echo.New()
	configure_web(web)

	slog.Info("setting up metrics")

	metrics := echo.New()
	configure_metrics(metrics)

	cleanup := []graceful.Cleanup{
		otelshutdown,
		metrics.Shutdown,
		web.Shutdown,
		db.DB().Shutdown,
	}

	graceful.Go(ctx, graceful.GrabAndGo(metrics.Start, ":"+PrometheusPort), errs)
	graceful.Go(ctx, graceful.GrabAndGo(web.Start, ":"+HTTPPort), errs)

	shared.Service().Banner()

	select {
	case err := <-errs:
		slog.Error("failed to start service", slog.Any("error", err.Error()))
	case <-sigterm:
		slog.Info("received shutdown signal")
	}

	code := graceful.Shutdown(ctx, cleanup, interrupt, 10*time.Second, 0)

	os.Exit(code)
}
