// Copyright © 2022, Breu Inc. <info@breu.io>. All rights reserved. 

package main

import (
	"log"
	"sync"

	"go.temporal.io/sdk/worker"

	"go.breu.io/ctrlplane/internal/db"
	"go.breu.io/ctrlplane/internal/providers/github"
	"go.breu.io/ctrlplane/internal/shared"
)

var wait sync.WaitGroup

func init() {
	shared.Service.ReadEnv()
	shared.Service.InitLogger()
	shared.EventStream.ReadEnv()
	shared.Temporal.ReadEnv()
	github.Github.ReadEnv()
	db.DB.ReadEnv()

	wait.Add(3)

	go func() {
		defer wait.Done()
		db.DB.InitSession()
	}()

	go func() {
		defer wait.Done()
		shared.EventStream.InitConnection()
	}()

	go func() {
		defer wait.Done()
		shared.Temporal.InitClient()
	}()

	wait.Wait()

	shared.Logger.Info("Initializing Service ... Done", "version", shared.Service.Version())
}

func main() {
	defer shared.Temporal.Client.Close()
	defer func() {
		if err := shared.Logger.Sync(); err != nil {
			panic(err)
		}
	}()

	queue := shared.Temporal.Queues[shared.ProvidersQueue].GetName()
	options := worker.Options{}
	wrkr := worker.New(shared.Temporal.Client, queue, options)

	workflows := &github.Workflows{}

	wrkr.RegisterWorkflow(workflows.OnInstall)
	wrkr.RegisterWorkflow(workflows.OnPush)
	wrkr.RegisterActivity(&github.Activities{})

	err := wrkr.Run(worker.InterruptCh())

	if err != nil {
		log.Fatal(err)
	}
}
