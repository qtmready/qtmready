package states

import (
	"errors"

	"github.com/google/uuid"
	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/core/repos/fns"
	"go.breu.io/quantm/internal/db/entities"
	"go.breu.io/quantm/internal/durable"
	"go.breu.io/quantm/internal/events"
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
)

type (
	// Repo defines the state for Repo Workflows. It embeds BaseState to inherit its functionality.
	Repo struct {
		*Base `json:"base"` // Base workflow state.

		Triggers BranchTriggers `json:"triggers"` // Branch triggers.
	}
)

// - signal handlers -

// OnPush is a signal handler for the push signal.
func (state *Repo) OnPush(ctx workflow.Context) durable.ChannelHandler {
	return func(rx workflow.ReceiveChannel, more bool) {
		event := &events.Event[eventsv1.RepoHook, eventsv1.Push]{}
		state.rx(ctx, rx, event)

		branch := fns.BranchNameFromRef(event.Payload.Ref)

		state.Triggers.add(branch, event.ID)
	}
}

// - query handlers -

// QueryBranchTrigger queries the parent branch for the specified branch.
func (state *Repo) QueryBranchTrigger(branch string) (uuid.UUID, error) {
	id, ok := state.Triggers.get(branch)
	if ok {
		return id, nil
	}

	return uuid.Nil, errors.New("branch not found")
}

// - state managers -

// RestartRecommended checks if the workflow should be continued as new.
func (state *Repo) RestartRecommended(ctx workflow.Context) bool {
	return workflow.GetInfo(ctx).GetContinueAsNewSuggested()
}

// NewRepo creates a new RepoState instance. It initializes BaseState using the provided context and
// hydrated repository data.
func NewRepo(repo *entities.Repo, msg *entities.Messaging) *Repo {
	base := &Base{Repo: repo, Messaging: msg}
	triggers := make(BranchTriggers)

	return &Repo{base, triggers} // Return new RepoState instance.
}
