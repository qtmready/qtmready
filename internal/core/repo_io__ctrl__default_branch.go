package core

import (
	"go.temporal.io/sdk/workflow"
)

func DefaultBranchCtrl(ctx workflow.Context, repo *Repo) error {
	return nil
}
