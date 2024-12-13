package git

import (
	"context"
	"log/slog"
)

type (
	// RepoOp represents the type of repository operation.
	RepoOp string
	// ResolveOp represents the type of resolve operation.
	ResolveOp string
	// CompareOp represents the type of comparison operation.
	CompareOp string

	// RepositoryError represents an error related to repository operations.
	RepositoryError struct {
		Op         RepoOp // Operation like "clone", "open"
		Repository *Repository
		internal   error
	}

	// ResolveError represents an error during revision/commit resolution.
	ResolveError struct {
		Op         ResolveOp // Operation like "resolve revision", "resolve commit"
		Ref        string    // Revision or commit reference
		Repository *Repository
		internal   error
	}

	// CompareError represents an error during comparison operations.
	CompareError struct {
		Op         CompareOp // Operation like "diff", "ancestor"
		From       string    // Source revision/commit
		To         string    // Target revision/commit
		Repository *Repository
		internal   error
	}
)

// - Repo Operation Constants -.
const (
	OpClone RepoOp = "clone"
	OpOpen  RepoOp = "open"
)

// - Resolve Operation Constants -.
const (
	OpResolveRevision ResolveOp = "resolve revision"
	OpResolveCommit   ResolveOp = "resolve commit"
)

// - Compare Operation Constants -.
const (
	OpDiff     CompareOp = "diff"
	OpAncestor CompareOp = "ancestor"
)

// - RepositoryError -

// Error method for RepositoryError.
func (e *RepositoryError) Error() string {
	return "repository error"
}

// Unwrap method for RepositoryError.
func (e *RepositoryError) Unwrap() error {
	return e.internal
}

// Wrap method to wrap the error.
func (e *RepositoryError) Wrap(err error) error {
	e.internal = err
	return e
}

func (e *RepositoryError) ReportError() error {
	return e.report(slog.LevelError)
}

func (e *RepositoryError) ReportWarn() error {
	return e.report(slog.LevelWarn)
}

func (e *RepositoryError) report(level slog.Level) error {
	attrs := []any{
		slog.String("operation", string(e.Op)),
		slog.String("repo_id", e.Repository.Entity.ID.String()),
		slog.String("repo_path", e.Repository.Path),
	}
	if e.internal != nil {
		attrs = append(attrs, slog.Any("details", e.internal))
	}

	slog.Log(context.Background(), level, e.Error(), attrs...)

	return e
}

// Helper function to create a new RepositoryError.
func NewRepositoryError(r *Repository, op RepoOp) *RepositoryError {
	return &RepositoryError{
		Op:         op,
		Repository: r,
	}
}

// - ResolveError -

// Error method for ResolveError.
func (e *ResolveError) Error() string {
	return "resolve error"
}

// Unwrap method for ResolveError.
func (e *ResolveError) Unwrap() error {
	return e.internal
}

// Wrap method to wrap the error.
func (e *ResolveError) Wrap(err error) error {
	e.internal = err
	return e
}

func (e *ResolveError) ReportError() error {
	return e.report(slog.LevelError)
}

func (e *ResolveError) ReportWarn() error {
	return e.report(slog.LevelWarn)
}

func (e *ResolveError) report(level slog.Level) error {
	attrs := []any{
		slog.String("operation", string(e.Op)),
		slog.String("repo_id", e.Repository.Entity.ID.String()),
		slog.String("repo_path", e.Repository.Path),
		slog.String("ref", e.Ref),
	}
	if e.internal != nil {
		attrs = append(attrs, slog.Any("details", e.internal))
	}

	slog.Log(context.Background(), level, e.Error(), attrs...)

	return e
}

// Helper function to create a new ResolveError.
func NewResolveError(r *Repository, op ResolveOp, ref string) *ResolveError {
	return &ResolveError{
		Op:         op,
		Ref:        ref,
		Repository: r,
	}
}

// - CompareError -

// Error method for CompareError.
func (e *CompareError) Error() string {
	return "compare error"
}

// Unwrap method for CompareError.
func (e *CompareError) Unwrap() error {
	return e.internal
}

// Wrap method to wrap the error.
func (e *CompareError) Wrap(err error) error {
	e.internal = err
	return e
}

func (e *CompareError) ReportError() error {
	return e.report(slog.LevelError)
}

func (e *CompareError) ReportWarn() error {
	return e.report(slog.LevelWarn)
}

func (e *CompareError) report(level slog.Level) error {
	attrs := []any{
		slog.String("operation", string(e.Op)),
		slog.String("repo_id", e.Repository.Entity.ID.String()),
		slog.String("repo_path", e.Repository.Path),
		slog.String("from", e.From),
		slog.String("to", e.To),
	}
	if e.internal != nil {
		attrs = append(attrs, slog.Any("details", e.internal))
	}

	slog.Log(context.Background(), level, e.Error(), attrs...)

	return e
}

// Helper function to create a new CompareError.
func NewCompareError(r *Repository, op CompareOp, from, to string) *CompareError {
	return &CompareError{
		Op:         op,
		From:       from,
		To:         to,
		Repository: r,
	}
}
