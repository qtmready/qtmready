package workflows

import (
	"log/slog"

	"go.breu.io/durex/dispatch"
	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/hooks/github/activities"
	"go.breu.io/quantm/internal/hooks/github/defs"
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
)

func Push(ctx workflow.Context, payload *defs.Push) error {
	acts := &activities.Push{}
	ctx = dispatch.WithDefaultActivityContext(ctx)

	var meta *defs.RepoEvent[eventsv1.RepoHook, eventsv1.Push]
	if err := workflow.
		ExecuteActivity(ctx, acts.ConvertToPushEvent, payload).
		Get(ctx, &meta); err != nil {
		return err
	}

	slog.Info("github/push: dispatching event ...", "repo", meta.Meta.Repo.Name, "", meta.Event)

	// TODO - need to confirm the signature
	return workflow.
		ExecuteActivity(ctx, acts.SignalCoreRepo, meta.Meta.Repo, defs.SignalWebhookPush, meta.Event).
		Get(ctx, nil)
}
