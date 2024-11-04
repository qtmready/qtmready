package githubwfs

import (
	"go.temporal.io/sdk/log"
	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/core/defs"
	githubdefs "go.breu.io/quantm/internal/hooks/github/defs"
)

type (
	PushWorkflowState struct {
		log log.Logger
	}
)

func Push(ctx workflow.Context) error {
	state := NewPushWorkflowState(ctx)
	selector := workflow.NewSelector(ctx)

	rqst := workflow.GetSignalChannel(ctx, githubdefs.SignalWebhookPush.String())
	selector.AddReceive(rqst, state.on_push(ctx))

	return nil
}

func (s *PushWorkflowState) on_push(ctx workflow.Context) defs.ChannelHandler {
	return func(rx workflow.ReceiveChannel, more bool) {}
}

func NewPushWorkflowState(ctx workflow.Context) *PushWorkflowState {
	return &PushWorkflowState{
		log: workflow.GetLogger(ctx),
	}
}
