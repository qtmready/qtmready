package code

import (
	"context"

	"go.breu.io/quantm/internal/core/defs" // Import Repo from defs
	"go.breu.io/quantm/internal/db"
)

type (
	// repoio represents the activities for the repo.
	repoio struct{}
)

// RepoIO creates and returns a new RepoIO object.
//
// Example:
//
//	repo_io := code.RepoIO()
func RepoIO() *repoio {
	return &repoio{}
}

// Get retrieves a repo from the database based on the provided parameters.
//
// Example:
//
//	repo, err := code.RepoIO().Get(ctx, db.QueryParams{"id": repo_id})
func (r *repoio) Get(ctx context.Context, params db.QueryParams) (*defs.Repo, error) {
	repo := &defs.Repo{}

	return repo, db.Get(repo, params)
}

// GetByID retrieves a repo from the database by their ID.
//
// Example:
//
//	repo, err := code.RepoIO().GetByID(ctx, repo_id)
func (r *repoio) GetByID(ctx context.Context, id string) (*defs.Repo, error) {
	repo := &defs.Repo{}

	return repo, db.Get(repo, db.QueryParams{"id": id})
}

// GetByCtrlID retrieves a repo from the database by their control ID.
//
// Example:
//
//	repo, err := code.RepoIO().GetByCtrlID(ctx, ctrl_id)
func (r *repoio) GetByCtrlID(ctx context.Context, ctrl_id string) (*defs.Repo, error) {
	repo := &defs.Repo{}

	return repo, db.Get(repo, db.QueryParams{"ctrl_id": ctrl_id})
}

// Save saves a repo to the database.
//
// Note: Always provide a pointer to the complete Repo object to avoid
// creating a copy. The Save method will update the provided Repo object
// with any changes made by the database.
//
// Example:
//
//	repo, err := code.RepoIO().Save(ctx, repo)
func (r *repoio) Save(ctx context.Context, repo *defs.Repo) (*defs.Repo, error) {
	return repo, db.Save(repo)
}
