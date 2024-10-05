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

	"go.breu.io/graceful"

	"go.breu.io/quantm/internal/core/defs"
	"go.breu.io/quantm/internal/core/kernel"
	"go.breu.io/quantm/internal/core/ws"
	"go.breu.io/quantm/internal/providers/github"
	"go.breu.io/quantm/internal/providers/slack"
	"go.breu.io/quantm/internal/shared"
	"go.breu.io/quantm/internal/shared/queues"
)

func main() {
	shared.Service().SetName("mothership")
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

	configure_core()
	configure_providers()
	configure_mutex()

	graceful.Go(ctx, ws.Queue().Start, rx_errors)

	cleanups := []graceful.Cleanup{
		queues.Core().Shutdown,
		queues.Providers().Shutdown,
		queues.Mutex().Shutdown,
		ws.Queue().Shutdown,
	}

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
