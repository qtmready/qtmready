package github

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/bradleyfalzon/ghinstallation/v2"
	gh "github.com/google/go-github/v62/github"

	"go.breu.io/quantm/internal/core"
	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/shared"
)

type (
	// RepoIO conforms to core.RepoIO interface.
	RepoIO struct{}
)

func (r *RepoIO) GetProviderInfo(ctx context.Context, id string) (*core.RepoIOProviderInfo, error) {
	repo := &Repo{}
	if err := db.Get(repo, db.QueryParams{"id": id}); err != nil {
		return nil, err
	}

	data := &core.RepoIOProviderInfo{
		RepoName:       repo.Name,
		DefaultBranch:  repo.DefaultBranch,
		ProviderID:     repo.GithubID.String(),
		RepoOwner:      strings.Split(repo.FullName, "/")[0],
		InstallationID: repo.InstallationID,
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

func (r *RepoIO) GetAllBranches(ctx context.Context, payload *core.RepoIOProviderInfo) ([]string, error) {
	branches := make([]string, 0)
	page := 1

	client, err := Instance().GetClientForInstallationID(payload.InstallationID)
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

	return branches, nil
}

// DetectChanges detects changes in a repository.
func (r *RepoIO) DetectChanges(ctx context.Context, payload *core.RepoIODetectChangesPayload) (*core.RepoIOChanges, error) {
	client, err := Instance().GetClientForInstallationID(payload.InstallationID)
	if err != nil {
		return nil, err
	}

	// TODO - move to some genernic function or activity
	// NOTE - need only repo URL... skip this call and get URL from payload or make it.
	repo, _, err := client.Repositories.Get(ctx, payload.RepoOwner, payload.RepoName)
	if err != nil {
		return nil, err
	}

	comparison, _, err := client.
		Repositories.
		CompareCommits(context.Background(), payload.RepoOwner, payload.RepoName, payload.DefaultBranch, payload.TargetBranch, nil)
	if err != nil {
		return nil, err
	}

	var changes, additions, deletions int

	files := make([]string, 0)

	for _, file := range comparison.Files {
		changes += file.GetChanges()
		additions += file.GetAdditions()
		deletions += file.GetDeletions()
		files = append(files, *file.Filename)
	}

	// detect changes struct
	dc := &core.RepoIOChanges{
		Added:      shared.Int64(additions),
		Removed:    shared.Int64(deletions),
		Modified:   files,
		Delta:      shared.Int64(changes),
		CompareUrl: comparison.GetHTMLURL(),
		RepoUrl:    repo.GetHTMLURL(),
	}

	return dc, nil
}

// Clone shallow clones a repository at a sepcific commit.
// see https://stackoverflow.com/a/76334845
func (r *RepoIO) TokenizedCloneURL(ctx context.Context, payload *core.RepoIOProviderInfo) (string, error) {
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
