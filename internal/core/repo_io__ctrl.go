package core

import (
	"go.temporal.io/sdk/workflow"
)

func _do(ctx workflow.Context, repo *Repo, branch, kind, action string, activity, payload, result any, keyvals ...any) error {
	logger := NewRepoIOWorkflowLogger(ctx, repo, kind, branch, action)
	logger.Info("init", keyvals...)

	if err := workflow.ExecuteActivity(ctx, activity, payload).Get(ctx, result); err != nil {
		logger.Warn("error", append(keyvals, "error", err)...)
		return err
	}

	logger.Info("success", keyvals...)

	return nil
}
