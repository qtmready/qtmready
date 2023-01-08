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
	"sync"

	"go.temporal.io/sdk/worker"

	"go.breu.io/ctrlplane/internal/db"
	"go.breu.io/ctrlplane/internal/providers/github"
	"go.breu.io/ctrlplane/internal/shared"
)

func init() {
	waitgroup := sync.WaitGroup{}

	shared.Service.ReadEnv()
	shared.Service.InitLogger()
	shared.EventStream.ReadEnv()
	shared.Temporal.ReadEnv()
	github.Github.ReadEnv()
	db.DB.ReadEnv()

	waitgroup.Add(3)
	shared.Logger.Info("initializing ...")

	go func() {
		defer waitgroup.Done()
		db.DB.InitSession()
	}()

	go func() {
		defer waitgroup.Done()
		shared.EventStream.InitConnection()
	}()

	go func() {
		defer waitgroup.Done()
		shared.Temporal.InitClient()
	}()

	waitgroup.Wait()

	shared.Logger.Info("initialized", "version", shared.Service.Version())
}

func main() {
	// graceful shutdown. see https://stackoverflow.com/a/46255965/228697.
	exitcode := 0
	defer func() { os.Exit(exitcode) }()
	defer func() { _ = shared.Logger.Sync() }()
	defer func() { _ = shared.EventStream.Drain() }()
	defer shared.Temporal.Client.Close()

	queue := shared.Temporal.Queues[shared.ProvidersQueue].GetName()
	options := worker.Options{}
	wrkr := worker.New(shared.Temporal.Client, queue, options)

	workflows := &github.Workflows{}

	wrkr.RegisterWorkflow(workflows.OnInstall)
	wrkr.RegisterWorkflow(workflows.OnInstallationRepositories)
	wrkr.RegisterWorkflow(workflows.OnPush)
	wrkr.RegisterWorkflow(workflows.OnPullRequest)

	wrkr.RegisterActivity(&github.Activities{})

	err := wrkr.Run(worker.InterruptCh())

	if err != nil {
		exitcode = 1
		return
	}
}
