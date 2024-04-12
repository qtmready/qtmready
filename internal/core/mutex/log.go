package mutex

import (
	"log/slog"

	"go.temporal.io/sdk/workflow"
)

func wfinfo(ctx workflow.Context, info *Info, msg string, extra ...any) {
	logger := workflow.GetLogger(ctx)
	id := info.ID
	caller := info.Caller.WorkflowExecution.ID

	extra = append(extra, slog.String("called by", caller), slog.String("resource id", id))

	logger.Info(msg, extra...)
}

func wferr(ctx workflow.Context, info *Info, msg string, err error) {
	logger := workflow.GetLogger(ctx)
	id := info.ID
	caller := info.Caller.WorkflowExecution.ID

	logger.Error(msg, slog.String("called by", caller), slog.String("resource id", id), slog.String("error", err.Error()))
}

func wfwarn(ctx workflow.Context, info *Info, msg string, err error) {
	logger := workflow.GetLogger(ctx)
	id := info.ID
	caller := info.Caller.WorkflowExecution.ID

	if err == nil {
		logger.Warn(msg, slog.String("called by", caller), slog.String("resource id", id))
	} else {
		logger.Warn(msg, slog.String("called by", caller), slog.String("resource id", id), slog.String("error", err.Error()))
	}
}
