package github

import (
	"context"

	gh "github.com/google/go-github/v62/github"

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

func (r *RepoIO) GetAllBranches(ctx context.Context, payload *core.RepoIOGetAllBranchesPayload) ([]string, error) {
	branches := make([]string, 0)
	page := 1

	client, err := Instance().GetClientForInstallation(payload.InstallationID)
	if err != nil {
		return branches, err
	}

	for {
		_branches, response, err := client.Repositories.ListBranches(
			ctx, payload.RepoOwner, payload.RepoName, &gh.BranchListOptions{
				ListOptions: gh.ListOptions{
					Page:    page,
					PerPage: 30, // Adjust this value as needed
				},
			},
		)

		if err != nil {
			return branches, err
		}

		for _, branch := range _branches {
			branches = append(branches, core.BranchNameFromRef(*branch.Name))
		}

		if response.NextPage == 0 {
			break
		}
	}

	// Get all branches for the repo
	return branches, nil
}
