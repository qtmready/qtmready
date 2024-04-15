package mutex

import (
	"log/slog"

	"go.temporal.io/sdk/workflow"
)

func wfdebug(ctx workflow.Context, info *Handler, msg string, extra ...any) {
	logger := workflow.GetLogger(ctx)
	id := info.ResourceID
	caller := info.Info.WorkflowExecution.ID

	extra = append(extra, slog.String("resource_id", id), slog.String("caller", caller))

	logger.Debug(msg, extra...)
}

func wfinfo(ctx workflow.Context, info *Handler, msg string, extra ...any) {
	logger := workflow.GetLogger(ctx)
	id := info.ResourceID
	caller := info.Info.WorkflowExecution.ID

	extra = append(extra, slog.String("resource_id", id), slog.String("caller", caller))

	logger.Info(msg, extra...)
}

func wferr(ctx workflow.Context, info *Handler, msg string, err error) {
	logger := workflow.GetLogger(ctx)
	id := info.ResourceID
	caller := info.Info.WorkflowExecution.ID

	logger.Error(msg, slog.String("resource_id", id), slog.String("caller", caller), slog.String("error", err.Error()))
}

func wfwarn(ctx workflow.Context, info *Handler, msg string, err error) {
	logger := workflow.GetLogger(ctx)
	id := info.ResourceID
	caller := info.Info.WorkflowExecution.ID

	if err == nil {
		logger.Warn(msg, slog.String("resource_id", id), slog.String("caller", caller))
	} else {
		logger.Warn(msg, slog.String("resource_id", id), slog.String("caller", caller), slog.String("error", err.Error()))
	}
}
