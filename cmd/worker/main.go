package main

import (
	"log"

	tworker "go.temporal.io/sdk/worker"

	"go.breu.io/ctrlplane/internal/conf"
	"go.breu.io/ctrlplane/internal/workflows"
)

func main() {
	defer conf.Temporal.Client.Close()
	options := tworker.Options{}
	worker := tworker.New(conf.Temporal.Client, conf.Temporal.Queues.Webhooks, options)

	worker.RegisterWorkflow(workflows.OnGithubInstall)

	err := worker.Run(tworker.InterruptCh())

	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	conf.InitService("webhooks-worker")
	conf.InitTemporal()
	conf.InitTemporalClient()
}
