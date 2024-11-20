package reposwfs

import (
	"go.temporal.io/sdk/workflow"

	reposdefs "go.breu.io/quantm/internal/core/repos/defs"
)

type (
	BaseState struct {
		*reposdefs.HypdratedRepo
	}
)

func NewBaseState(ctx workflow.Context, hydrated *reposdefs.HypdratedRepo) *BaseState {
	return &BaseState{hydrated}
}
