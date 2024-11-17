package reposwfs

import (
	"go.temporal.io/sdk/workflow"

	reposdefs "go.breu.io/quantm/internal/core/repos/defs"
	reposstate "go.breu.io/quantm/internal/core/repos/state"
)

// Repo manages the event loop for a repository, acting as a central router to orchestrate repository workflows.
func Repo(ctx workflow.Context, repo *reposdefs.FullRepo) error {
	selector := workflow.NewSelector(ctx)

	// TODO - need to discuss how to set the state for repo and base state
	state := &reposstate.RepoState{}

	// push event
	push := workflow.GetSignalChannel(ctx, reposdefs.RepoIOSignalPush.String())
	selector.AddReceive(push, state.OnPush(ctx))

	return nil
}
