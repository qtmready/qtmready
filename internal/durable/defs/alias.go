package defs

import (
	"go.temporal.io/sdk/workflow"
)

type (
	FutureHandler    func(workflow.Future)               // FutureHandler is the signature of the future handler for temporal.
	ChannelHandler   func(workflow.ReceiveChannel, bool) // ChannelHandler is the signature of the channel handler for temporal.
	CoroutineHandler func(workflow.Context)              // CoroutineHandler is the signature of the coroutine handler for temporal.
)
