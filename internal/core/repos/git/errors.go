package git

import (
	"fmt"
)

type (
	// RepositoryError represents an error related to repository operations.
	RepositoryError struct {
		Op         RepoOp // Operation like "clone", "open"
		Path       string
		Repo       string
		Repository *Repository
		internal   error
	}

	// RepoOp represents the type of repository operation.
	RepoOp string

	// ResolveError represents an error during revision/commit resolution.
	ResolveError struct {
		Op         ResolveOp // Operation like "resolve revision", "resolve commit"
		Ref        string    // Revision or commit reference
		Repository *Repository
		internal   error
	}

	// ResolveOp represents the type of resolve operation.
	ResolveOp string

	// CompareError represents an error during comparison operations.
	CompareError struct {
		Op         CompareOp // Operation like "diff", "ancestor"
		From       string    // Source revision/commit
		To         string    // Target revision/commit
		Repository *Repository
		internal   error
	}

	// CompareOp represents the type of comparison operation.
	CompareOp string
)

// Repo Operation Constants.
const (
	OpClone RepoOp = "clone"
	OpOpen  RepoOp = "open"
)

// Resolve Operation Constants.
const (
	OpResolveRevision ResolveOp = "resolve revision"
	OpResolveCommit   ResolveOp = "resolve commit"
)

// Compare Operation Constants.
const (
	OpDiff     CompareOp = "diff"
	OpAncestor CompareOp = "ancestor"
)

// Error method for RepositoryError.
func (e *RepositoryError) Error() string {
	if e.internal != nil {
		return fmt.Sprintf("repository %s error: path: %s, repo: %s, details: %v", e.Op, e.Path, e.Repo, e.internal)
	}

	return fmt.Sprintf("repository %s error: path: %s, repo: %s, internal: %s", e.Op, e.Path, e.Repo, e.internal.Error())
}

// Unwrap method for RepositoryError.
func (e *RepositoryError) Unwrap() error { return e.internal }

// Error method for ResolveError.
func (e *ResolveError) Error() string {
	if e.internal != nil {
		return fmt.Sprintf("resolve %s error: ref: %s, details: %v", e.Op, e.Ref, e.internal)
	}

	return fmt.Sprintf("resolve %s error: ref: %s", e.Op, e.Ref)
}

// Unwrap method for ResolveError.
func (e *ResolveError) Unwrap() error { return e.internal }

// Error method for CompareError.
func (e *CompareError) Error() string {
	if e.internal != nil {
		return fmt.Sprintf("compare %s error: from: %s, to: %s, details: %v", e.Op, e.From, e.To, e.internal)
	}

	return fmt.Sprintf("compare %s error: from: %s, to: %s", e.Op, e.From, e.To)
}

// Unwrap method for CompareError.
func (e *CompareError) Unwrap() error { return e.internal }

// Helper function to create a new RepositoryError.
func NewRepositoryError(r *Repository, op RepoOp) *RepositoryError {
	return &RepositoryError{
		Op:         op,
		Path:       r.Path,
		Repo:       r.Entity.ID.String(),
		Repository: r,
	}
}

// Helper function to create a new ResolveError.
func NewResolveError(r *Repository, op ResolveOp, ref string) *ResolveError {
	return &ResolveError{
		Op:         op,
		Ref:        ref,
		Repository: r,
	}
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

// Wrap method to wrap the error.
func (e *RepositoryError) Wrap(err error) error {
	e.internal = err
	return e
}

// Wrap method to wrap the error.
func (e *ResolveError) Wrap(err error) error {
	e.internal = err
	return e
}

// Wrap method to wrap the error.
func (e *CompareError) Wrap(err error) error {
	e.internal = err
	return e
}
