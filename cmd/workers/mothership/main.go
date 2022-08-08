package main

import (
	"log"
	"sync"

	"go.temporal.io/sdk/worker"

	"go.breu.io/ctrlplane/internal/cmn"
	"go.breu.io/ctrlplane/internal/db"
	"go.breu.io/ctrlplane/internal/integrations/github"
)

var wait sync.WaitGroup

func init() {
	cmn.Service.ReadEnv()
	cmn.Service.InitLogger()
	cmn.EventStream.ReadEnv()
	cmn.Temporal.ReadEnv()
	github.Github.ReadEnv()
	db.DB.ReadEnv()

	wait.Add(3)

	go func() {
		defer wait.Done()
		db.DB.InitSession()
	}()

	go func() {
		defer wait.Done()
		cmn.EventStream.InitConnection()
	}()

	go func() {
		defer wait.Done()
		cmn.Temporal.InitClient()
	}()

	wait.Wait()

	cmn.Log.Info("Initializing Service ... Done")
}

func main() {
	defer cmn.Temporal.Client.Close()

	queue := cmn.Temporal.Queues.Integrations
	options := worker.Options{}
	wrkr := worker.New(cmn.Temporal.Client, queue, options)

	workflows := &github.Workflows{}

	wrkr.RegisterWorkflow(workflows.OnInstall)
	wrkr.RegisterWorkflow(workflows.OnPush)
	wrkr.RegisterActivity(&github.Activity{})

	err := wrkr.Run(worker.InterruptCh())

	if err != nil {
		log.Fatal(err)
	}
}
