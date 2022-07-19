package main

import (
	"log"
	"sync"

	"go.temporal.io/sdk/worker"

	"go.breu.io/ctrlplane/internal/common"
	"go.breu.io/ctrlplane/internal/db"
	"go.breu.io/ctrlplane/internal/integrations/github"
)

var wait sync.WaitGroup

func init() {
	common.Service.ReadEnv()
	common.Service.InitLogger()
	common.EventStream.ReadEnv()
	common.Temporal.ReadEnv()
	github.Github.ReadEnv()
	db.DB.ReadEnv()

	wait.Add(3)

	go func() {
		defer wait.Done()
		db.DB.InitSession()
	}()

	go func() {
		defer wait.Done()
		common.EventStream.InitConnection()
	}()

	go func() {
		defer wait.Done()
		common.Temporal.InitClient()
	}()

	wait.Wait()

	common.Logger.Info("Initializing Service ... Done")
}

func main() {
	defer common.Temporal.Client.Close()

	queue := common.Temporal.Queues.Integrations
	options := worker.Options{}
	wrkr := worker.New(common.Temporal.Client, queue, options)

	workflows := &github.Workflows{}

	wrkr.RegisterWorkflow(workflows.OnInstall)
	wrkr.RegisterWorkflow(workflows.OnPush)
	wrkr.RegisterActivity(&github.Activity{})

	err := wrkr.Run(worker.InterruptCh())

	if err != nil {
		log.Fatal(err)
	}
}
