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
)

func main() {
	ctx := context.Background()
	interrupt := make(chan any, 1) // channel to signal the shutdown to goroutines.
	errs := make(chan error, 1)    // create a channel to listen to errors.

	// Handle termination signals (SIGINT, SIGTERM, SIGQUIT)
	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, os.Interrupt)

	// init service
	shared.Service().SetName("mothership")
	shared.Logger().Info(
		"starting ...",
		slog.Any("service", shared.Service().GetName()),
		slog.String("version", shared.Service().GetVersion()),
	)

	kernel.Instance(
		kernel.WithRepoProvider(defs.RepoProviderGithub, &github.RepoIO{}),
		kernel.WithMessageProvider(defs.MessageProviderSlack, &slack.Activities{}),
	)

	cleanups := []graceful.Cleanup{}

	graceful.Go(ctx, graceful.FreezeAndFizzle(q_core, interrupt), errs)
	graceful.Go(ctx, graceful.FreezeAndFizzle(q_provider, interrupt), errs)

	shared.Service().Banner()

	select {
	case <-terminate:
		shared.Logger().Info("received shutdown signal")
	case err := <-errs:
		shared.Logger().Error("unable to start ", "error", err)
	}

	code := graceful.Shutdown(ctx, cleanups, interrupt, 10*time.Second, 0)

	os.Exit(code)
}

func q_core(interrupt <-chan any) error {
	worker := shared.Temporal().Queue(shared.CoreQueue).Worker(shared.Temporal().Client())

	worker.RegisterWorkflow(code.RepoCtrl)
	worker.RegisterWorkflow(code.TrunkCtrl)
	worker.RegisterWorkflow(code.BranchCtrl)
	worker.RegisterWorkflow(code.QueueCtrl)

	worker.RegisterActivity(&code.Activities{})

	// TODO: this will not work if we have more than one provider.
	worker.RegisterActivity(&github.RepoIO{})
	worker.RegisterActivity(&slack.Activities{})

	worker.RegisterActivity(mutex.PrepareMutexActivity)

	return worker.Run(interrupt)
}

func q_provider(interrupt <-chan any) error {
	worker := shared.Temporal().Queue(shared.ProvidersQueue).Worker(shared.Temporal().Client())

	github_workflows := &github.Workflows{}

	worker.RegisterWorkflow(github_workflows.OnInstallationEvent)
	worker.RegisterWorkflow(github_workflows.OnInstallationRepositoriesEvent)
	worker.RegisterWorkflow(github_workflows.PostInstall)
	worker.RegisterWorkflow(github_workflows.OnPushEvent)
	worker.RegisterWorkflow(github_workflows.OnCreateOrDeleteEvent)
	worker.RegisterWorkflow(github_workflows.OnPullRequestEvent)
	worker.RegisterWorkflow(github_workflows.OnWorkflowRunEvent)

	worker.RegisterActivity(&github.Activities{})

	return worker.Run(interrupt)
}
