package states

import (
	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/db/entities"
)

type (
	Branch struct {
		*Base  `json:"base"` // Base workflow state.
		Branch string        `json:"branch"`
	}
)

func (state *Branch) Init(ctx workflow.Context) {
	state.Base.Init(ctx)
}

func NewBranch(repo *entities.Repo, msg *entities.Messaging, branch string) *Branch {
	base := &Base{Repo: repo, Messaging: msg}

	return &Branch{base, branch}
}
