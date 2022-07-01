package main

import (
	"log"

	tw "go.temporal.io/sdk/worker"

	"go.breu.io/ctrlplane/internal/conf"
	"go.breu.io/ctrlplane/internal/temporal/activities"
	"go.breu.io/ctrlplane/internal/temporal/workflows"
)

func main() {
	defer conf.Temporal.Client.Close()

	queue := conf.Temporal.Queues.Webhooks
	options := tw.Options{}
	worker := tw.New(conf.Temporal.Client, queue, options)

	worker.RegisterWorkflow(workflows.OnGithubInstall)
	worker.RegisterActivity(activities.GetOrCreateGithubInstallation)

	err := worker.Run(tw.InterruptCh())

	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	conf.InitService("worker::webhooks")
	conf.InitTemporal()
	conf.InitTemporalClient()
}
