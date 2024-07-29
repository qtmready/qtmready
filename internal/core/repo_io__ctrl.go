package core

import (
	"go.temporal.io/sdk/workflow"
)

// _do is a helper function that executes an activity within a workflow context, logging the
// activity execution and handling any errors that may occur.
//
// It takes the following parameters:
//
//   - ctx: the workflow context with activity execution options already set.
//   - repo: the repository object
//   - branch: the branch name
//   - kind: the kind of operation being performed
//   - action: the action being performed
//   - activity: the activity function to execute
//   - payload: the input data for the activity
//   - result: a pointer to store the activity result
//   - keyvals: additional key-value pairs to include in the log
//
// If the activity execution is successful, it logs the success. If there is an error, it logs the error and returns it.
func _do(ctx workflow.Context, repo *Repo, branch, kind, action string, activity, payload, result any, keyvals ...any) error {
	logger := NewRepoIOWorkflowLogger(ctx, repo, kind, branch, action)
	logger.Info("init", keyvals...)

	if err := workflow.ExecuteActivity(ctx, activity, payload).Get(ctx, result); err != nil {
		logger.Warn("error", append(keyvals, "error", err)...)
		return err
	}

	logger.Info("result", append(keyvals, "result", result)...)
	logger.Info("success", keyvals...)

	return nil
}
