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

	"go.breu.io/quantm/internal/core/code"
	"go.breu.io/quantm/internal/core/defs"
	"go.breu.io/quantm/internal/core/kernel"
	"go.breu.io/quantm/internal/core/mutex"
	"go.breu.io/quantm/internal/providers/github"
	"go.breu.io/quantm/internal/providers/slack"
	"go.breu.io/quantm/internal/shared"
	"go.breu.io/quantm/internal/shared/graceful"
	"go.breu.io/quantm/internal/shared/queue"
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

	shared.Logger().Info("starting workers...")

	core_queue := shared.Temporal().Queue(shared.CoreQueue)
	configure_core(core_queue)

	provider_queue := shared.Temporal().Queue(shared.ProvidersQueue)
	configure_provider(provider_queue)

	cleanups := []graceful.Cleanup{}

	graceful.Go(ctx, graceful.WrapRelease(core_queue.Listen, release), rx_errors)
	graceful.Go(ctx, graceful.WrapRelease(provider_queue.Listen, release), rx_errors)

	shared.Service().Banner()

	select {
	case rx := <-terminate:
		slog.Info("main: received shutdown signal, attempting graceful shutdown ...", "signal", rx.String())
	case err := <-rx_errors:
		slog.Error("main: unable to start ...", "error", err.Error())
	}

	code := graceful.Shutdown(ctx, cleanups, release, timeout, 0)
	if code == 1 {
		slog.Warn("main: failed to shutdown gracefully, exiting ...")
	}

	os.Exit(code)
}

func configure_core(q queue.Queue) {
	worker := q.Worker(shared.Temporal().Client())

	worker.RegisterWorkflow(code.RepoCtrl)
	worker.RegisterWorkflow(code.TrunkCtrl)
	worker.RegisterWorkflow(code.BranchCtrl)
	worker.RegisterWorkflow(code.QueueCtrl)

	worker.RegisterActivity(&code.Activities{})

	// TODO: this will not work if we have more than one provider.
	worker.RegisterActivity(&github.RepoIO{})
	worker.RegisterActivity(&slack.Activities{})

	worker.RegisterActivity(mutex.PrepareMutexActivity)
}

func configure_provider(q queue.Queue) {
	worker := q.Worker(shared.Temporal().Client())

	github_workflows := &github.Workflows{}

	worker.RegisterWorkflow(github_workflows.OnInstallationEvent)
	worker.RegisterWorkflow(github_workflows.OnInstallationRepositoriesEvent)
	worker.RegisterWorkflow(github_workflows.PostInstall)
	worker.RegisterWorkflow(github_workflows.OnPushEvent)
	worker.RegisterWorkflow(github_workflows.OnCreateOrDeleteEvent)
	worker.RegisterWorkflow(github_workflows.OnPullRequestEvent)
	worker.RegisterWorkflow(github_workflows.OnWorkflowRunEvent)

	worker.RegisterActivity(&github.Activities{})
}
