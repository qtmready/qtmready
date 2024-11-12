package githubwfs

import (
	"go.temporal.io/sdk/log"
	"go.temporal.io/sdk/workflow"

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

	var result *githubdefs.Trigger[commonv1.RepoHook, eventsv1.Push]
	if err := workflow.
		ExecuteActivity(ctx, acts.ConvertToPushEvent, payload).
		Get(ctx, &result); err != nil {
		return err
	}

	return nil
}
