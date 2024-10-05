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
	"go.breu.io/graceful"

	"go.breu.io/quantm/internal/core/defs"
	"go.breu.io/quantm/internal/core/kernel"
	"go.breu.io/quantm/internal/core/ws"
	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/providers/github"
	"go.breu.io/quantm/internal/providers/slack"
	"go.breu.io/quantm/internal/shared"
)

const (
	HTTPPort       = "8000"
	PrometheusPort = "9090"
)

func main() {
	shared.Service().SetName("api")
	shared.Logger().Info("main: init ...", "service", shared.Service().GetName(), "version", shared.Service().GetVersion())

	ctx := context.Background()
	release := make(chan any, 1)
	rx_errors := make(chan error)
	timeout := time.Second * 10

	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, os.Interrupt)

	kernel.Instance(
		kernel.WithRepoProvider(defs.RepoProviderGithub, &github.RepoIO{}),
		kernel.WithMessageProvider(defs.MessageProviderSlack, &slack.Activities{}),
	)

	otelshutdown, err := observe(ctx, shared.Service().GetName(), shared.Service().GetVersion())
	if err != nil {
		slog.Error("main: failed to setup opentelemetry, exiting ...", slog.Any("error", err.Error()))
		rx_errors <- err

		return
	}

	hub := ws.Instance()
	hub.SetAuthFn(user_id)

	api := echo.New()
	configure_api(api)

	metrics := echo.New()
	configure_metrics(metrics)

	cleanups := []graceful.Cleanup{
		otelshutdown,
		api.Shutdown,
		metrics.Shutdown,
		db.DB().Shutdown,
		hub.Stop,
	}

	graceful.Go(ctx, graceful.GrabAndGo(metrics.Start, ":"+PrometheusPort), rx_errors)
	graceful.Go(ctx, graceful.GrabAndGo(api.Start, ":"+HTTPPort), rx_errors)

	shared.Service().Banner()

	select {
	case rx := <-terminate:
		slog.Info("main: shutdown requested ...", "signal", rx.String())
	case err := <-rx_errors:
		slog.Error("main: unable to start ...", "error", err.Error())
	}

	code := graceful.Shutdown(ctx, cleanups, release, timeout, 0)

	if code == 1 {
		slog.Warn("main: failed to shutdown gracefully, exiting ...")
	} else {
		slog.Info("main: shutdown complete, exiting ...")
	}

	os.Exit(code)
}
