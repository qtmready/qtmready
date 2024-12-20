package states

import (
	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/db/entities"
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
)

type (
	Trunk struct {
		*Base
		MergeQueue *Sequencer[int64, eventsv1.MergeQueue]
		IsEmpty    bool
	}
)

// - queue process -.
func (state Trunk) ProcessItem() {}

func (state *Trunk) Init(ctx workflow.Context) {
	state.Base.Init(ctx)
	state.MergeQueue.Init(ctx)
}

func NewTrunk(repo *entities.Repo, chat *entities.ChatLink) *Trunk {
	return &Trunk{
		&Base{Repo: repo, ChatLink: chat},
		NewSequencer[int64, eventsv1.MergeQueue](),
		true,
	}
}
