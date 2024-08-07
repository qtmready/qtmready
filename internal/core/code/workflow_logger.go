package code

import (
	"go.temporal.io/sdk/log"
	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/core/defs"
)

type (
	RepoIOWorkflowLogger struct {
		repo   *defs.Repo
		kind   string
		branch string
		action string
		logger log.Logger
	}

	LogWriter func(msg string, keyvals ...any)
)

func NewRepoIOWorkflowLogger(ctx workflow.Context, repo *defs.Repo, kind, branch, action string) *RepoIOWorkflowLogger {
	logger := workflow.GetLogger(ctx)

	return &RepoIOWorkflowLogger{repo, kind, branch, action, logger}
}

func (r *RepoIOWorkflowLogger) Info(msg string, keyvals ...any) {
	r.write(r.logger.Info, msg, keyvals...)
}

func (r *RepoIOWorkflowLogger) Warn(msg string, keyvals ...any) {
	r.write(r.logger.Warn, msg, keyvals...)
}

func (r *RepoIOWorkflowLogger) Error(msg string, keyvals ...any) {
	r.write(r.logger.Error, msg, keyvals...)
}

func (r *RepoIOWorkflowLogger) Debug(msg string, keyvals ...any) {
	r.write(r.logger.Debug, msg, keyvals...)
}

func (r *RepoIOWorkflowLogger) prefix() string {
	prefix := r.kind

	if r.branch != "" {
		prefix += "/" + r.branch
	}

	if r.action != "" {
		prefix += "/" + r.action
	}

	return prefix + ": "
}

func (r *RepoIOWorkflowLogger) write(writer LogWriter, msg string, keyvals ...any) {
	keyvals = append(keyvals, "repo_id", r.repo.ID.String())
	keyvals = append(keyvals, "provider", r.repo.Provider.String())
	keyvals = append(keyvals, "provider_id", r.repo.ProviderID)
	keyvals = append(keyvals, "default_branch", r.repo.DefaultBranch)

	if r.branch != "" {
		keyvals = append(keyvals, "branch", r.branch)
	}

	if r.action != "" {
		keyvals = append(keyvals, "scope", r.action)
	}

	writer(r.prefix()+msg, keyvals...)
}
