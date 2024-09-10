// Copyright © 2023, Breu, Inc. <info@breu.io>
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
	"os"
	"os/signal"
	"syscall"

	"go.breu.io/quantm/internal/core/code"
	"go.breu.io/quantm/internal/core/defs"
	"go.breu.io/quantm/internal/core/kernel"
	"go.breu.io/quantm/internal/core/mutex"
	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/providers/github"
	"go.breu.io/quantm/internal/providers/slack"
	"go.breu.io/quantm/internal/shared"
)

func main() {
	shared.Service().SetName("mothership")
	// graceful shutdown. see https://stackoverflow.com/a/46255965/228697.
	exitcode := 0
	defer func() { os.Exit(exitcode) }()
	defer shared.Temporal().Client().Close()
	defer db.DB().Session.Close()

	providerWrkr := shared.Temporal().Worker(shared.ProvidersQueue)
	coreWrkr := shared.Temporal().Worker(shared.CoreQueue)

	kernel.Instance(
		kernel.WithRepoProvider(defs.RepoProviderGithub, &github.RepoIO{}),
		kernel.WithMessageProvider(defs.MessageProviderSlack, &slack.Activities{}),
	)

	githubwfs := &github.Workflows{}

	// provider workflows
	providerWrkr.RegisterWorkflow(githubwfs.OnInstallationEvent)
	providerWrkr.RegisterWorkflow(githubwfs.OnInstallationRepositoriesEvent)
	providerWrkr.RegisterWorkflow(githubwfs.PostInstall)
	providerWrkr.RegisterWorkflow(githubwfs.OnPushEvent)
	providerWrkr.RegisterWorkflow(githubwfs.OnCreateOrDeleteEvent)
	providerWrkr.RegisterWorkflow(githubwfs.OnPullRequestEvent)
	providerWrkr.RegisterWorkflow(githubwfs.OnWorkflowRunEvent)

	// provider activities
	providerWrkr.RegisterActivity(&github.Activities{})
	providerWrkr.RegisterActivity(&slack.Activities{})

	// mutex workflow
	coreWrkr.RegisterWorkflow(mutex.MutexWorkflow)
	providerWrkr.RegisterWorkflow(mutex.MutexWorkflow)

	// code workflows
	coreWrkr.RegisterWorkflow(code.RepoCtrl)
	coreWrkr.RegisterWorkflow(code.TrunkCtrl)
	coreWrkr.RegisterWorkflow(code.BranchCtrl)
	coreWrkr.RegisterWorkflow(code.QueueCtrl)

	// core activities
	coreWrkr.RegisterActivity(&code.Activities{})

	// RepoIO & MessageIO
	coreWrkr.RegisterActivity(&github.RepoIO{})
	coreWrkr.RegisterActivity(&slack.Activities{})

	// mutex activity
	coreWrkr.RegisterActivity(mutex.PrepareMutexActivity)
	providerWrkr.RegisterActivity(mutex.PrepareMutexActivity)

	// start worker for provider queue
	err := providerWrkr.Start()
	if err != nil {
		exitcode = 1
		return
	}

	// start worker for core queue
	err = coreWrkr.Start()
	if err != nil {
		exitcode = 1
		return
	}

	shared.Service().Banner()

	quit := make(chan os.Signal, 1)                      // create a channel to listen to quit signals.
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // setting up the signals to listen to.
	<-quit                                               // wait for quit signal.

	shared.Logger().Info("Shutting down workers...")

	providerWrkr.Stop()
	coreWrkr.Stop()

	shared.Logger().Info("Workers stopped. Exiting...")

	exitcode = 1
}
