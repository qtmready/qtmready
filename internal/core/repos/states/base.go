package states

import (
	"fmt"

	"go.temporal.io/sdk/log"
	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/db/entities"
)

type (
	// Base represents the base state for repository workflows.  It encapsulates
	// core data and provides logging utilities.
	Base struct {
		Repo      *entities.Repo      `json:"repo"`      // Repository entity.
		Messaging *entities.Messaging `json:"messaging"` // Messaging entity.

		logger log.Logger // Workflow logger.
	}
)

// rx wraps workflow.ReceiveChannel.Receive, adding logging.  It receives a message
// from the specified Temporal channel. The target parameter must be a pointer to the
// data structure expected to be received.
func (s *Base) rx(ctx workflow.Context, ch workflow.ReceiveChannel, target any) {
	s.logger.Info(fmt.Sprintf("rx: %s", ch.Name()))
	ch.Receive(ctx, target)
}

// dispatch wraps activity execution with logging.
func (s *Base) dispatch(ctx workflow.Context, action string) {}

// Init initializes the base state with the provided context.
func (s *Base) Init(ctx workflow.Context) {
	s.logger = workflow.GetLogger(ctx)
}
