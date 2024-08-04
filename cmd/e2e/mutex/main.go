package main

import (
	"context"
	"crypto/rand"
	"log/slog"
	"math/big"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/core/mutex"
	"go.breu.io/quantm/internal/shared"
	"go.breu.io/quantm/internal/shared/queue"
)

type (
	Data = map[uuid.UUID]time.Duration
)

var (
	id, _ = uuid.NewV7()
)

func main() {
	shared.Logger().Info("starting ...")
	wk := configure(shared.CoreQueue)

	if err := wk.Start(); err != nil {
		slog.Error("Unable to start worker", slog.String("error", err.Error()))
	}

	opts := shared.Temporal().Queue(shared.CoreQueue).WorkflowOptions(
		shared.WithWorkflowBlock("parent"),
		shared.WithWorkflowBlockID(id.String()),
	)

	if _, err := shared.Temporal().Client().ExecuteWorkflow(context.Background(), opts, ParentWorkflow); err != nil {
		slog.Error("Unable to execute workflow", slog.String("error", err.Error()))
	}

	quit := make(chan os.Signal, 1)                      // create a channel to listen to quit signals.
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // setting up the signals to listen to.
	<-quit
}

func configure(queue queue.Name) worker.Worker {
	worker := shared.Temporal().Worker(queue)
	worker.RegisterWorkflow(mutex.MutexWorkflow)
	worker.RegisterWorkflow(ParentWorkflow)
	worker.RegisterWorkflow(ChildWorkflow)

	worker.RegisterActivity(mutex.PrepareMutexActivity)

	return worker
}

func ParentWorkflow(ctx workflow.Context) error {
	logger := workflow.GetLogger(ctx)
	queue := make(Data, 0)
	futures := make([]workflow.Future, 0)

	for range 50 {
		workflow.SideEffect(ctx, func(workflow.Context) any {
			n, _ := rand.Int(rand.Reader, big.NewInt(30))
			wait := time.Duration(n.Int64()) * time.Second
			id, _ := uuid.NewV7()
			queue[id] = wait

			return nil
		})
	}

	for id := range queue {
		opts := shared.Temporal().Queue(shared.CoreQueue).ChildWorkflowOptions(
			shared.WithWorkflowParent(ctx),
			shared.WithWorkflowBlock("child"),
			shared.WithWorkflowBlockID(id.String()),
		)
		childctx := workflow.WithChildOptions(ctx, opts)
		future := workflow.ExecuteChildWorkflow(childctx, ChildWorkflow, id, queue[id])

		futures = append(futures, future)
	}

	for _, future := range futures {
		logger.Info("waiting for child workflows to complete")

		_ = future.Get(ctx, nil)
	}

	return nil
}

func ChildWorkflow(ctx workflow.Context, id uuid.UUID, timeout time.Duration) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("Starting child workflow", slog.String("id", id.String()), slog.String("timeout", timeout.String()))

	lock := mutex.New(
		ctx,
		mutex.WithResourceID("repo.xyz"),
		mutex.WithTimeout(timeout+(10*time.Second)),
	)

	// Prepare the lock means that get the reference to running Mutex workflow and schedule a new lock on it. If there is no Mutex workflow
	// running, then start a new Mutex workflow and schedule a lock on it.
	if err := lock.Prepare(ctx); err != nil {
		return err // or handle error
	}

	// Acquire acquires the lock. If we do not handle the error.
	if err := lock.Acquire(ctx); err != nil {
		return err // or handle error
	}

	// Do so work here.
	if err := workflow.Sleep(ctx, timeout); err != nil {
		return err // or handle error
	}

	// Release the lock.
	if err := lock.Release(ctx); err != nil {
		return err // or handle error
	}

	// Cleanup tries to shutdown the Mutex workflow if there are no more locks waiting.
	if err := lock.Cleanup(ctx); err != nil {
		return err // or handle error
	}

	return nil
}
