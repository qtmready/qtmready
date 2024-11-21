package reposwfs

import (
	"go.temporal.io/sdk/workflow"

	reposdefs "go.breu.io/quantm/internal/core/repos/defs"
	"go.breu.io/quantm/internal/db/entities"
	"go.breu.io/quantm/internal/durable"
	"go.breu.io/quantm/internal/events"
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
)

type (
	// RepoState defines the state for Repo Workflows. It embeds BaseState to inherit its functionality.
	RepoState struct {
		*BaseState `json:"base"` // Base workflow state.

		Triggers BranchTriggers `json:"triggers"` // Branch triggers.
	}
)

// Repo manages the event loop for a repository. It acts as a central router, orchestrating repository workflows.
//
// Repo uses Temporal's workflow selector to concurrently handle signals. The function initializes a RepoState and
// registers a signal handler for the reposdefs.RepoIOSignalPush signal.  Currently, the signal handler is a stub.
// Temporal workflow context is passed as input, along with hydrated repository data.  The function returns an error
// if one occurs during workflow execution; otherwise it returns nil.
//
// FIXME: start the function with the state instead of HydratedRepo.If we start from a state, It will be easier to
// restart the workflow as new.
func Repo(ctx workflow.Context, state *RepoState) error {
	state.init(ctx)

	selector := workflow.NewSelector(ctx)

	push := workflow.GetSignalChannel(ctx, reposdefs.SignalPush.String())
	selector.AddReceive(push, state.on_push(ctx))

	for !state.refresh_urged(ctx) {
		selector.Select(ctx)
	}

	return workflow.NewContinueAsNewError(ctx, Repo, state)
}

// - signal handlers -

// on_push is a signal handler for the push signal.
func (state *RepoState) on_push(ctx workflow.Context) durable.ChannelHandler {
	return func(rx workflow.ReceiveChannel, more bool) {
		event := &events.Event[eventsv1.RepoHook, eventsv1.Push]{}
		state.rx(ctx, rx, event)

		state.Triggers.add("branch", event.ID) // TODO: add the right branch.
	}
}

// - state managers -

func (state *RepoState) refresh_urged(ctx workflow.Context) bool {
	return workflow.GetInfo(ctx).GetContinueAsNewSuggested()
}

// NewRepoState creates a new RepoState instance. It initializes BaseState using the provided context and
// hydrated repository data.
func NewRepoState(repo *entities.Repo, msg *entities.Messaging) *RepoState {
	base := &BaseState{Repo: repo, Messaging: msg}
	triggers := make(BranchTriggers)

	return &RepoState{base, triggers} // Return new RepoState instance.
}
