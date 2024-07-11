package timers

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

// Now returns the current time in the context of the workflow.
// This is a side effect that ensures that time is deterministic across replays.
func Now(ctx workflow.Context) time.Time {
	var now time.Time

	_ = workflow.SideEffect(ctx, func(_ctx workflow.Context) any { return time.Now() }).Get(&now)

	return now
}
