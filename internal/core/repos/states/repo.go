package states

import (
	"errors"

	"github.com/google/uuid"
	"go.breu.io/durex/dispatch"
	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/core/repos/activities"
	"go.breu.io/quantm/internal/core/repos/defs"
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

		do *activities.Repo
	}
)

// - local -

func (state *Repo) signal_branch(ctx workflow.Context, branch string, event any) error {
	ctx = dispatch.WithDefaultActivityContext(ctx)

	next := NewBranch(state.Repo, state.Messaging, branch)
	payload := &activities.SignalBranchPayload{Signal: defs.SignalPush, Repo: state.Repo, Branch: branch}

	return workflow.ExecuteActivity(ctx, state.do.SignalBranch, payload, event, next).Get(ctx, nil)
}

// - signal handlers -

// OnPush is a signal handler for the push signal.
func (state *Repo) OnPush(ctx workflow.Context) durable.ChannelHandler {
	return func(rx workflow.ReceiveChannel, more bool) {
		event := &events.Event[eventsv1.RepoHook, eventsv1.Push]{}
		state.rx(ctx, rx, event)

		branch := fns.BranchNameFromRef(event.Payload.Ref)

		state.Triggers.add(branch, event.ID)

		if err := state.signal_branch(ctx, branch, event); err != nil {
			state.logger.Warn("unable to signal branch", "repo", state.Repo.ID, "branch", branch, "error", err.Error())
		}
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

func (state *Repo) Init(ctx workflow.Context) {
	state.Base.Init(ctx)
}

// NewRepo creates a new RepoState instance. It initializes BaseState using the provided context and
// hydrated repository data.
func NewRepo(repo *entities.Repo, msg *entities.Messaging) *Repo {
	base := &Base{Repo: repo, Messaging: msg}
	triggers := make(BranchTriggers)

	return &Repo{base, triggers, &activities.Repo{}} // Return new RepoState instance.
}
