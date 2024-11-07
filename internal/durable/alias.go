package durable

import (
	"go.breu.io/quantm/internal/durable/defs"
)

type (

	// FutureHandler is the signature for the function that handles the future result of an operation for temporal workflows.
	FutureHandler = defs.FutureHandler

	// ChannelHandler is the signature for the function that handles messages received on a channel for temporal workflows.
	ChannelHandler = defs.ChannelHandler

	// CoroutineHandler is the signature for the function that handles the results of a goroutine for temporal workflows.
	CoroutineHandler = defs.CoroutineHandler
)
