package githubwfs

import (
	"log/slog"

	"go.breu.io/durex/dispatch"
	"go.temporal.io/sdk/workflow"

	githubacts "go.breu.io/quantm/internal/hooks/github/activities"
	githubdefs "go.breu.io/quantm/internal/hooks/github/defs"
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
)

func Push(ctx workflow.Context, payload *githubdefs.Push) error {
	acts := &githubacts.Push{}
	ctx = dispatch.WithDefaultActivityContext(ctx)

	var meta *githubdefs.RepoEvent[eventsv1.RepoHook, eventsv1.Push]
	if err := workflow.
		ExecuteActivity(ctx, acts.ConvertToPushEvent, payload).
		Get(ctx, &meta); err != nil {
		return err
	}

	slog.Info("github/push: dispatching event ...", "repo", meta.Meta.Repo.Name, "", meta.Event)

	// TODO - need to confirm the signature
	return workflow.
		ExecuteActivity(ctx, acts.SignalCoreRepo, meta.Meta.Repo, githubdefs.SignalWebhookPush, meta.Event).
		Get(ctx, nil)
}
