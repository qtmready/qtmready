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

	"github.com/sourcegraph/conc"
	"go.temporal.io/sdk/worker"

	"go.breu.io/ctrlplane/internal/core"
	"go.breu.io/ctrlplane/internal/db"
	"go.breu.io/ctrlplane/internal/providers/github"
	"go.breu.io/ctrlplane/internal/shared"
)

func init() {
	waitgroup := conc.NewWaitGroup()
	defer waitgroup.Wait()

	shared.Service.ReadEnv()
	shared.Service.InitLogger(3)
	shared.EventStream.ReadEnv()
	shared.Temporal.ReadEnv()
	github.Github.ReadEnv()
	db.DB.ReadEnv()

	// shared.Temporal.ServerHost = "127.0.0.1"
	// db.DB.Hosts = append(db.DB.Hosts, "127.0.0.1")
	waitgroup.Go(db.DB.InitSession)
	// waitgroup.Go(shared.EventStream.InitConnection)
	waitgroup.Go(shared.Temporal.InitClient)

	shared.Logger.Info("initialized", "version", shared.Service.Version())

	core.Core.Init()
}

func main() {
	// graceful shutdown. see https://stackoverflow.com/a/46255965/228697.
	exitcode := 0
	defer func() { os.Exit(exitcode) }()
	defer func() { _ = shared.Logger.Sync() }()
	// defer func() { _ = shared.EventStream.Drain() }()
	defer shared.Temporal.Client.Close()

	core.Core.ProvidersMap[core.RepoProviderGithub] = github.Github

	queue := shared.Temporal.Queues[shared.ProvidersQueue].GetName()
	coreQueue := shared.Temporal.Queues[shared.CoreQueue].GetName()

	options := worker.Options{OnFatalError: func(err error) { shared.Logger.Error("Fatal error during worker execution", err) }}
	wrkr := worker.New(shared.Temporal.Client, queue, options)
	coreWrkr := worker.New(shared.Temporal.Client, coreQueue, options)

	ghwfs := &github.Workflows{}
	cwfs := &core.Workflows{}

	// provider workflows
	wrkr.RegisterWorkflow(ghwfs.OnInstallationEvent)
	wrkr.RegisterWorkflow(ghwfs.OnInstallationRepositoriesEvent)
	wrkr.RegisterWorkflow(ghwfs.OnPushEvent)
	wrkr.RegisterWorkflow(ghwfs.OnPullRequestEvent)

	// provider activities
	wrkr.RegisterActivity(&github.Activities{})
	wrkr.RegisterActivity(github.Github.GetLatestCommitforRepo)

	// core workflows
	coreWrkr.RegisterWorkflow(cwfs.OnPullRequestWorkflow)
	coreWrkr.RegisterWorkflow(cwfs.MutexWorkflow)
	coreWrkr.RegisterWorkflow(cwfs.DeploymentWorkflow)
	coreWrkr.RegisterWorkflow(cwfs.GetAssetsWorkflow)
	coreWrkr.RegisterWorkflow(cwfs.ProvisionInfraWorkflow)
	coreWrkr.RegisterWorkflow(cwfs.DeProvisionInfraWorkflow)

	// core activities
	coreWrkr.RegisterActivity(&core.Activities{})

	// start worker for provider queue
	err := wrkr.Start()
	if err != nil {
		exitcode = 1
		return
	}

	defer wrkr.Stop()

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

	shared.Logger.Info("Exiting....")

	exitcode = 1
}
