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

	"go.temporal.io/sdk/worker"

	"go.breu.io/ctrlplane/internal/core"
	"go.breu.io/ctrlplane/internal/core/mutex"
	"go.breu.io/ctrlplane/internal/core/resources/gcp"
	"go.breu.io/ctrlplane/internal/db"
	"go.breu.io/ctrlplane/internal/providers/github"
	"go.breu.io/ctrlplane/internal/shared"
)

func main() {
	// graceful shutdown. see https://stackoverflow.com/a/46255965/228697.
	exitcode := 0
	defer func() { os.Exit(exitcode) }()
	defer func() { _ = shared.Logger().Sync() }()
	defer shared.Temporal().Client().Close()
	defer db.DB().Session.Close()

	core.Instance(
		core.WithRepoProvider(core.RepoProviderGithub, &github.Activities{}),
		core.WithCloudResource(core.CloudProviderGCP, core.DriverCloudrun, &gcp.CloudRunConstructor{}),
	)

	providerQueue := shared.Temporal().Queue(shared.ProvidersQueue).Name()
	coreQueue := shared.Temporal().Queue(shared.CoreQueue).Name()

	options := worker.Options{OnFatalError: func(err error) { shared.Logger().Error("Fatal error during worker execution", err) }}
	providerWrkr := worker.New(shared.Temporal().Client(), providerQueue, options)
	coreWrkr := worker.New(shared.Temporal().Client(), coreQueue, options)

	ghwfs := &github.Workflows{}
	cwfs := &core.Workflows{}
	gcpwfs := gcp.Workflows{}

	// provider workflows
	providerWrkr.RegisterWorkflow(ghwfs.OnInstallationEvent)
	providerWrkr.RegisterWorkflow(ghwfs.OnInstallationRepositoriesEvent)
	providerWrkr.RegisterWorkflow(ghwfs.OnPushEvent)
	providerWrkr.RegisterWorkflow(ghwfs.OnPullRequestEvent)

	// provider activities
	providerWrkr.RegisterActivity(&github.Activities{})

	// mutex workflow
	coreWrkr.RegisterWorkflow(mutex.Workflow)

	// core workflows
	coreWrkr.RegisterWorkflow(cwfs.StackController)
	coreWrkr.RegisterWorkflow(cwfs.Deploy)
	coreWrkr.RegisterWorkflow(cwfs.GetAssets)
	coreWrkr.RegisterWorkflow(cwfs.ProvisionInfra)
	coreWrkr.RegisterWorkflow(cwfs.DeProvisionInfra)
	coreWrkr.RegisterWorkflow(gcpwfs.DeployCloudRun)
	coreWrkr.RegisterWorkflow(gcpwfs.UpdateTraffic)

	// core activities
	coreWrkr.RegisterActivity(&core.Activities{})
	coreWrkr.RegisterActivity(&gcp.Activities{})

	// start worker for provider queue
	err := providerWrkr.Start()
	if err != nil {
		exitcode = 1
		return
	}

	defer providerWrkr.Stop()

	// start worker for core queue
	err = coreWrkr.Start()
	if err != nil {
		exitcode = 1
		return
	}

	defer coreWrkr.Stop()

	quit := make(chan os.Signal, 1)                      // create a channel to listen to quit signals.
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // setting up the signals to listen to.
	<-quit                                               // wait for quit signal.

	shared.Logger().Info("Exiting....")

	exitcode = 1
}
