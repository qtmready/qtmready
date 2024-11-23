package workflows

import (
	"errors"

	"github.com/google/uuid"
	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/core/repos/defs"
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
func Repo(ctx workflow.Context, state *RepoState) error {
	state.init(ctx)

	selector := workflow.NewSelector(ctx)

	// - query handlers -
	if err := workflow.SetQueryHandler(ctx, defs.QueryRepoForEventParent.String(), state.q__branch_trigger); err != nil {
		return err
	}

	// - signal handlers -

	push := workflow.GetSignalChannel(ctx, defs.SignalPush.String())
	selector.AddReceive(push, state.on_push(ctx))

	// - event loop -

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

// - query handlers -

// q__branch_trigger queries the parent branch for the specified branch.
func (state *RepoState) q__branch_trigger(branch string) (uuid.UUID, error) {
	id, ok := state.Triggers.get(branch)
	if ok {
		return id, nil
	}

	return uuid.Nil, errors.New("branch not found")
}

// - state managers -

// refresh_urged checks if the workflow should be continued as new.
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
