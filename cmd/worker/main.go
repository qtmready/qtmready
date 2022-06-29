package main

import (
	"log"

	_tworker "go.temporal.io/sdk/worker"

	"go.breu.io/ctrlplane/internal/conf"
	"go.breu.io/ctrlplane/internal/workflows"
)

func main() {
	defer conf.Temporal.Client.Close()
	worker := _tworker.New(conf.Temporal.Client, conf.Temporal.Queues.Webhooks, _tworker.Options{})
	worker.RegisterWorkflow(workflows.OnGithubInstall)

	err := worker.Run(_tworker.InterruptCh())

	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	conf.InitService("webhooks-worker")
	conf.InitTemporal()
	conf.InitTemporalClient()
}
