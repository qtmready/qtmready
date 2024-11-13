package reposstate

import (
	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/durable/defs"
)

type (
	// RepoState defines the state for Repo Workflows.
	RepoState struct{}
)

func (state *RepoState) OnPush(ctx workflow.Context) defs.ChannelHandler {
	return func(rx workflow.ReceiveChannel, more bool) {

	}
}
