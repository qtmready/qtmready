package main

import (
	"context"
	"log"

	_tmprlactivity "go.temporal.io/sdk/activity"
	_tmprlworker "go.temporal.io/sdk/worker"
	_tmprlworkflow "go.temporal.io/sdk/workflow"

	"go.breu.io/ctrlplane/internal/conf"
	"go.breu.io/ctrlplane/internal/workflows"
)

func main() {
	defer conf.Temporal.Client.Close()

	worker := _tmprlworker.New(conf.Temporal.Client, conf.Temporal.Queues.Webhooks, _tmprlworker.Options{})
	worker.RegisterWorkflow(workflows.OnGithubInstall)
	// worker.RegisterActivity(firstActivity)

	err := worker.Run(_tmprlworker.InterruptCh())

	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	conf.InitService("webhooks-worker")
	conf.InitTemporal()
	conf.InitTemporalClient()
}

func firstWorkflow(ctx _tmprlworkflow.Context, input string) error {
	logger := _tmprlworkflow.GetLogger(ctx)
	logger.Info("first workflow started")
	return nil
}

func firstActivity(ctx context.Context, input string) error {
	logger := _tmprlactivity.GetLogger(ctx)
	logger.Info("first activity started with input: {}", input)
	return nil
}
