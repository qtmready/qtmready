package states

import (
	"fmt"

	"go.breu.io/durex/dispatch"
	"go.temporal.io/sdk/log"
	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/db/entities"
)

type (
	// Base represents the base state for repository workflows.  It encapsulates
	// core data and provides logging utilities.
	Base struct {
		Repo *entities.Repo      `json:"repo"` // Repository entity.
		Chat *entities.Messaging `json:"chat"` // Messaging entity.

		logger log.Logger // Workflow logger.
	}
)

// - private

// rx wraps workflow.ReceiveChannel.Receive, adding logging.  It receives a message
// from the specified Temporal channel. The target parameter must be a pointer to the
// data structure expected to be received.
func (state *Base) rx(ctx workflow.Context, ch workflow.ReceiveChannel, target any) {
	state.logger.Info(fmt.Sprintf("rx: %s", ch.Name()))
	ch.Receive(ctx, target)
}

// run wraps workflow.ExecuteActivity with logging with the default activity context. If you need to
// provide a custom context, use run_ex.
func (state *Base) run(ctx workflow.Context, action string, activity, event, result any, keyvals ...any) error {
	state.logger.Info(fmt.Sprintf("dispatch(%s): init ...", action), keyvals...)

	ctx = dispatch.WithDefaultActivityContext(ctx)

	if err := workflow.ExecuteActivity(ctx, activity, event).Get(ctx, result); err != nil {
		state.logger.Error(fmt.Sprintf("dispatch(%s): error", action), keyvals...)
		return err
	}

	state.logger.Info(fmt.Sprintf("dispatch(%s): success", action), keyvals...)

	return nil
}

// - public

// RestartRecommended checks if the workflow should be continued as new.
func (state *Base) RestartRecommended(ctx workflow.Context) bool {
	return workflow.GetInfo(ctx).GetContinueAsNewSuggested()
}

// Init initializes the base state with the provided context.
func (state *Base) Init(ctx workflow.Context) {
	state.logger = workflow.GetLogger(ctx)
}
