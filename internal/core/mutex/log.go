package mutex

import (
	"log/slog"

	"go.temporal.io/sdk/workflow"
)

func wfdebug(ctx workflow.Context, info *Info, msg string, extra ...any) {
	logger := workflow.GetLogger(ctx)
	id := info.ResourceID
	caller := info.Caller.WorkflowExecution.ID

	extra = append(extra, slog.String("called_by", caller), slog.String("resource_id", id))

	logger.Debug(msg, extra...)
}

func wfinfo(ctx workflow.Context, info *Info, msg string, extra ...any) {
	logger := workflow.GetLogger(ctx)
	id := info.ResourceID
	caller := info.Caller.WorkflowExecution.ID

	extra = append(extra, slog.String("called_by", caller), slog.String("resource_id", id))

	logger.Info(msg, extra...)
}

func wferr(ctx workflow.Context, info *Info, msg string, err error) {
	logger := workflow.GetLogger(ctx)
	id := info.ResourceID
	caller := info.Caller.WorkflowExecution.ID

	logger.Error(msg, slog.String("called_by", caller), slog.String("resource_id", id), slog.String("error", err.Error()))
}

func wfwarn(ctx workflow.Context, info *Info, msg string, err error) {
	logger := workflow.GetLogger(ctx)
	id := info.ResourceID
	caller := info.Caller.WorkflowExecution.ID

	if err == nil {
		logger.Warn(msg, slog.String("called_by", caller), slog.String("resource_id", id))
	} else {
		logger.Warn(msg, slog.String("called_by", caller), slog.String("resource_id", id), slog.String("error", err.Error()))
	}
}
