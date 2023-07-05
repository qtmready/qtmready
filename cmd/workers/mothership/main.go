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

	"go.breu.io/ctrlplane/internal/core"
	"go.breu.io/ctrlplane/internal/core/mutex"
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
	)

	/**
	 * Providers
	 **/
	providerwrkr := shared.Temporal().
		Queue(shared.ProvidersQueue).
		Worker(shared.Temporal().Client())
	ghwrkflos := &github.Workflows{}

	providerwrkr.RegisterWorkflow(ghwrkflos.OnInstallationEvent)
	providerwrkr.RegisterWorkflow(ghwrkflos.OnInstallationRepositoriesEvent)
	providerwrkr.RegisterWorkflow(ghwrkflos.OnPushEvent)
	providerwrkr.RegisterWorkflow(ghwrkflos.OnPullRequestEvent)

	providerwrkr.RegisterActivity(&github.Activities{})

	if err := providerwrkr.Start(); err != nil {
		exitcode = 1
		return
	}
	defer providerwrkr.Stop()

	/**
	 * Core
	 **/

	corewrkr := shared.Temporal().
		Queue(shared.CoreQueue).
		Worker(shared.Temporal().Client())
	corewrkflos := &core.Workflows{}

	corewrkr.RegisterWorkflow(mutex.Workflow)
	corewrkr.RegisterWorkflow(corewrkflos.StackController)
	corewrkr.RegisterWorkflow(corewrkflos.Deploy)
	corewrkr.RegisterWorkflow(corewrkflos.GetAssets)
	corewrkr.RegisterWorkflow(corewrkflos.ProvisionInfra)
	corewrkr.RegisterWorkflow(corewrkflos.DeProvisionInfra)

	corewrkr.RegisterActivity(&core.Activities{})

	if err := corewrkr.Start(); err != nil {
		exitcode = 1
		return
	}
	defer corewrkr.Stop()

	/**
	 * Worker successfully started. Announcing ...
	 **/
	shared.Service().Banner()

	/**
	 * Graceful Shutdown
	  **/

	quit := make(chan os.Signal, 1)                      // create a channel to listen to quit signals.
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // setting up the signals to listen to.
	<-quit                                               // wait for quit signal.

	shared.Logger().Info("Exiting....")

	exitcode = 1
}
