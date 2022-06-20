package main

import (
	"context"
	"log"

	_activity "go.temporal.io/sdk/activity"
	_client "go.temporal.io/sdk/client"
	_worker "go.temporal.io/sdk/worker"
	_workflow "go.temporal.io/sdk/workflow"
)

func main() {
	client, err := _client.Dial(_client.Options{})
	if err != nil {
		log.Fatal(err)
	}

	defer client.Close()

	worker := _worker.New(client, "hello-world", _worker.Options{})
	worker.RegisterWorkflow(firstWorkflow)
	worker.RegisterActivity(firstActivity)

	err = worker.Run(_worker.InterruptCh())

	if err != nil {
		log.Fatal(err)
	}
}

func firstWorkflow(ctx _workflow.Context, input string) error {
	logger := _workflow.GetLogger(ctx)
	logger.Info("first workflow started")
	return nil
}

func firstActivity(ctx context.Context, input string) error {
	logger := _activity.GetLogger(ctx)
	logger.Info("first activity started with input: {}", input)
	return nil
}
