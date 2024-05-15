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
	"os"
	"os/signal"
	"syscall"

	"go.breu.io/quantm/internal/core"
	"go.breu.io/quantm/internal/core/mutex"
	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/providers/gcp/cloudrun"
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

	core.Instance(
		core.WithRepoProvider(core.RepoProviderGithub, &github.Activities{}),
		core.WithCloudResource(core.CloudProviderGCP, core.DriverCloudrun, &cloudrun.Constructor{}),
		core.WithMessageProvider(core.MessageProviderSlack, &slack.Activities{}),
	)

	githubwfs := &github.Workflows{}
	stackwfs := &core.StackWorkflows{}
	repowfs := &core.RepoWorkflows{}

	// provider workflows
	providerWrkr.RegisterWorkflow(githubwfs.OnInstallationEvent)
	providerWrkr.RegisterWorkflow(githubwfs.OnInstallationRepositoriesEvent)
	providerWrkr.RegisterWorkflow(githubwfs.RefreshDefaultBranch)
	providerWrkr.RegisterWorkflow(githubwfs.OnPushEvent)
	providerWrkr.RegisterWorkflow(githubwfs.OnPullRequestEvent)
	providerWrkr.RegisterWorkflow(githubwfs.OnLabelEvent)
	providerWrkr.RegisterWorkflow(githubwfs.OnWorkflowRunEvent)

	// provider activities
	providerWrkr.RegisterActivity(&github.Activities{})
	providerWrkr.RegisterActivity(&slack.Activities{})

	// mutex workflow
	coreWrkr.RegisterWorkflow(mutex.Workflow)
	providerWrkr.RegisterWorkflow(mutex.Workflow)

	// stack workflows
	coreWrkr.RegisterWorkflow(stackwfs.StackController)
	coreWrkr.RegisterWorkflow(stackwfs.Deploy)
	coreWrkr.RegisterWorkflow(stackwfs.GetAssets)
	coreWrkr.RegisterWorkflow(stackwfs.ProvisionInfra)
	coreWrkr.RegisterWorkflow(stackwfs.DeProvisionInfra)

	// repo workflows
	coreWrkr.RegisterWorkflow(repowfs.BranchController)
	coreWrkr.RegisterWorkflow(repowfs.StaleBranchDetection)
	coreWrkr.RegisterWorkflow(repowfs.PollMergeQueue)

	// core activities
	coreWrkr.RegisterActivity(&core.Activities{})
	coreWrkr.RegisterActivity(&cloudrun.Activities{})

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
