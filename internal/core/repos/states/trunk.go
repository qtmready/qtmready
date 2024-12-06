package states

import (
	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/db/entities"
)

type (
	Trunk struct {
		*Base
	}
)

func (state *Trunk) Init(ctx workflow.Context) {
	state.Base.Init(ctx)
}

func NewTrunk(repo *entities.Repo, msg *entities.ChatLink) *Trunk {
	return &Trunk{
		&Base{Repo: repo, ChatLink: msg},
	}
}
