package reposwfs

import (
	"go.temporal.io/sdk/workflow"
)

type (
	TrunkState struct {
		*BaseState
	}
)

func Trunk(ctx workflow.Context) error { return nil }
