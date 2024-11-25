package states

import (
	"go.temporal.io/sdk/workflow"

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
	}
)

func (state *Branch) calculate_complexity() {}

func (state *Branch) OnPush(ctx workflow.Context) durable.ChannelHandler {
	return func(ch workflow.ReceiveChannel, more bool) {
		event := &events.Event[eventsv1.RepoHook, eventsv1.Push]{}
		state.rx(ctx, ch, event)

		state.LatestCommit = fns.GetLatestCommit(event.Payload)
	}
}

func (state *Branch) Init(ctx workflow.Context) {
	state.Base.Init(ctx)
}

func NewBranch(repo *entities.Repo, msg *entities.Messaging, branch string) *Branch {
	base := &Base{Repo: repo, Messaging: msg}

	return &Branch{Base: base, Branch: branch}
}
