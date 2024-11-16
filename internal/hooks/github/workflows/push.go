package githubwfs

import (
	"log/slog"

	"go.breu.io/durex/dispatch"
	"go.temporal.io/sdk/log"
	"go.temporal.io/sdk/workflow"

	githubacts "go.breu.io/quantm/internal/hooks/github/activities"
	githubdefs "go.breu.io/quantm/internal/hooks/github/defs"
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
)

type (
	PushWorkflowState struct {
		log log.Logger
	}
)

func Push(ctx workflow.Context, payload *githubdefs.Push) error {
	acts := &githubacts.Push{}
	ctx = dispatch.WithDefaultActivityContext(ctx)

	var eventory *githubdefs.RepoEvent[eventsv1.RepoHook, eventsv1.Push]
	if err := workflow.
		ExecuteActivity(ctx, acts.ConvertToPushEvent, payload).
		Get(ctx, &eventory); err != nil {
		return err
	}

	slog.Info("github/push: dispatching event ...")

	// TODO - need to confirm the signature
	return workflow.ExecuteActivity(
		ctx, acts.SignalCoreRepo, eventory.Repo, githubdefs.SignalWebhookPush, eventory.Event).
		Get(ctx, nil)
}
