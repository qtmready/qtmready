package reposwfs

import (
	"fmt"

	"go.temporal.io/sdk/log"
	"go.temporal.io/sdk/workflow"

	reposdefs "go.breu.io/quantm/internal/core/repos/defs"
	"go.breu.io/quantm/internal/db/entities"
)

type (
	// BaseState represents the base state for repository workflows.  It encapsulates
	// core data and provides logging utilities.
	BaseState struct {
		Repo      *entities.Repo      `json:"repo"`      // Repository entity.
		Messaging *entities.Messaging `json:"messaging"` // Messaging entity.

		log log.Logger // Workflow logger.
	}
)

// rx wraps workflow.ReceiveChannel.Receive, adding logging.  It receives a message
// from the specified Temporal channel. The target parameter must be a pointer to the
// data structure expected to be received.
func (state *BaseState) rx(ctx workflow.Context, ch workflow.ReceiveChannel, target any) {
	state.log.Info(fmt.Sprintf("rx: %s", ch.Name()))
	ch.Receive(ctx, target)
}

// dispatch wraps activity execution with logging.
func (state *BaseState) dispatch(ctx workflow.Context, action string) {}

// NewBaseState creates a new BaseState instance. It retrieves a workflow logger
// from the provided context and initializes the BaseState with repository and
// messaging data.
func NewBaseState(ctx workflow.Context, hydrated *reposdefs.HypdratedRepo) *BaseState {
	logger := workflow.GetLogger(ctx)
	return &BaseState{Repo: hydrated.Repo, Messaging: hydrated.Messaging, log: logger}
}
