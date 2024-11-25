package states

import (
	"github.com/google/uuid"
	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/core/repos/activities"
	"go.breu.io/quantm/internal/core/repos/fns"
	"go.breu.io/quantm/internal/db/entities"
	"go.breu.io/quantm/internal/durable"
	"go.breu.io/quantm/internal/events"
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
)

type (
	Branch struct {
		*Base        `json:"base"`    // Base workflow state.
		Branch       string           `json:"branch"`
		LatestCommit *eventsv1.Commit `json:"latest_commit"`

		acts *activities.Branch
	}
)

// OnPush is a channel handler for the push event.
func (state *Branch) OnPush(ctx workflow.Context) durable.ChannelHandler {
	return func(ch workflow.ReceiveChannel, more bool) {
		event := &events.Event[eventsv1.RepoHook, eventsv1.Push]{}
		state.rx(ctx, ch, event)
		state.calculate_complexity(ctx, event)
	}
}

// calculate_complexity calculates the complexity of the change against trunk.
func (state *Branch) calculate_complexity(ctx workflow.Context, event *events.Event[eventsv1.RepoHook, eventsv1.Push]) {
	state.LatestCommit = fns.GetLatestCommit(event.Payload)

	path := ""

	_ = workflow.SideEffect(ctx, func(ctx workflow.Context) any {
		return uuid.New().String()
	}).Get(&path)

	payload := &activities.ClonePayload{Repo: state.Repo, Branch: state.Branch, Hook: event.Context.Hook, Path: path}

	_ = state.run(ctx, "clone", state.acts.Clone, payload, nil)
}

// Init initializes the branch state with the provided context.
func (state *Branch) Init(ctx workflow.Context) {
	state.Base.Init(ctx)
}

func NewBranch(repo *entities.Repo, msg *entities.Messaging, branch string) *Branch {
	base := &Base{Repo: repo, Messaging: msg}

	return &Branch{Base: base, Branch: branch, acts: &activities.Branch{}}
}
