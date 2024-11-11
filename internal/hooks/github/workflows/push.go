package githubwfs

import (
	"go.temporal.io/sdk/log"
	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/events"
	githubacts "go.breu.io/quantm/internal/hooks/github/activities"
	githubdefs "go.breu.io/quantm/internal/hooks/github/defs"
	commonv1 "go.breu.io/quantm/internal/proto/ctrlplane/common/v1"
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
)

type (
	PushWorkflowState struct {
		log log.Logger
	}
)

func Push(ctx workflow.Context, payload *githubdefs.Push) error {
	acts := &githubacts.Push{}

	var event *events.Event[commonv1.RepoHook, eventsv1.Push]

	if err := workflow.
		ExecuteActivity(ctx, acts.ConvertToPushEvent, payload).
		Get(ctx, &event); err != nil {
		return err
	}

	return nil
}
