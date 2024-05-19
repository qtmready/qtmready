package github

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bradleyfalzon/ghinstallation/v2"
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

func (r *RepoIO) GetAllBranches(ctx context.Context, payload *core.RepoIOInfoPayload) ([]string, error) {
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

// DetectChanges detects changes in a repository.
func (r *RepoIO) DetectChanges(ctx context.Context, payload *core.RepoSignalPushPayload) (*core.RepoIOChanges, error) {
	return nil, nil
}

// Clone shallow clones a repository at a sepcific commit.
// see https://stackoverflow.com/a/76334845
func (r *RepoIO) TokenizedCloneURL(ctx context.Context, payload *core.RepoIOInfoPayload) (string, error) {
	installation, err := ghinstallation.New(
		http.DefaultTransport, Instance().AppID, payload.InstallationID.Int64(), []byte(Instance().PrivateKey),
	)

	if err != nil {
		return "", err
	}

	token, err := installation.Token(ctx)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("https://git:%s@github.com/%s/%s.git", token, payload.RepoOwner, payload.RepoName), nil
}
