package github

import (
	"context"

	"go.breu.io/quantm/internal/core"
	"go.breu.io/quantm/internal/db"
)

type (
	// RepoIO conforms to core.RepoIO interface.
	RepoIO struct{}
)

func (r *RepoIO) GetRepoData(ctx context.Context, id string) (*core.RepoIORepoData, error) {
	repo := &Repo{}
	if err := db.Get(repo, db.QueryParams{"id": id}); err != nil {
		return nil, err
	}

	data := &core.RepoIORepoData{
		Name:          repo.Name,
		DefaultBranch: repo.DefaultBranch,
		ProviderID:    repo.GithubID.String(),
	}

	return data, nil
}

func (r *RepoIO) SetEarlyWarning(ctx context.Context, id string, value bool) error {
	repo := &Repo{}
	if err := db.Get(repo, db.QueryParams{"id": id}); err != nil {
		return err
	}

	repo.HasEarlyWarning = value

	if err := db.Update(repo); err != nil {
		return err
	}

	return nil
}

func (r *RepoIO) GetAllBranches(ctx context.Context) error {
	return nil
}
