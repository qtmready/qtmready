package reposwfs

import (
	"go.temporal.io/sdk/workflow"

	reposdefs "go.breu.io/quantm/internal/core/repos/defs"
	"go.breu.io/quantm/internal/durable"
)

type (
	// RepoState defines the state for Repo Workflows.
	RepoState struct {
		*BaseState `json:"base"`
	}
)

// Repo manages the event loop for a repository, acting as a central router to orchestrate repository workflows.
func Repo(ctx workflow.Context, hydrated_repo *reposdefs.HypdratedRepo) error {
	selector := workflow.NewSelector(ctx)

	// TODO - need to discuss how to set the state for repo and base state
	state := NewRepoState(ctx, hydrated_repo)

	// push event
	push := workflow.GetSignalChannel(ctx, reposdefs.RepoIOSignalPush.String())
	selector.AddReceive(push, state.OnPush(ctx))

	return nil
}

func (state *RepoState) OnPush(ctx workflow.Context) durable.ChannelHandler {
	return func(rx workflow.ReceiveChannel, more bool) {}
}

func NewRepoState(ctx workflow.Context, hydrated *reposdefs.HypdratedRepo) *RepoState {
	base := NewBaseState(ctx, hydrated)
	return &RepoState{base}
}
