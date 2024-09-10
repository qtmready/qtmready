// Copyright © 2023, Breu, Inc. <info@breu.io>
//
// We hereby irrevocably grant you an additional license to use the Software under the Apache License, Version 2.0 that
// is effective on the second anniversary of the date we make the Software available. On or after that date, you may use
// the Software under the Apache License, Version 2.0, in which case the following will apply:
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
// the License.
//
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
// an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.

package github

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/bradleyfalzon/ghinstallation/v2"
	gh "github.com/google/go-github/v62/github"

	"go.breu.io/quantm/internal/core/code"
	"go.breu.io/quantm/internal/core/defs"
	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/shared"
)

type (
	RepoIO struct{}
)

func (r *RepoIO) GetProviderInfo(ctx context.Context, id string) (*defs.RepoIOProviderInfo, error) {
	repo := &Repo{}
	if err := db.Get(repo, db.QueryParams{"id": id}); err != nil {
		return nil, err
	}

	data := &defs.RepoIOProviderInfo{
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

func (r *RepoIO) GetAllBranches(ctx context.Context, payload *defs.RepoIOProviderInfo) ([]string, error) {
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
			branches = append(branches, code.BranchNameFromRef(*branch.Name))
		}

		if response.NextPage == 0 {
			break
		}
	}

	return branches, nil
}

// DetectChanges detects changes in a repository.
func (r *RepoIO) DetectChanges(ctx context.Context, payload *defs.RepoIODetectChangesPayload) (*defs.RepoIOChanges, error) {
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
	dc := &defs.RepoIOChanges{
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
func (r *RepoIO) TokenizedCloneURL(ctx context.Context, payload *defs.RepoIOProviderInfo) (string, error) {
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

// MergePR Branch in default repo branch.
// TODO - need to refine.
// NOTE - to optimze the logic need to make a logic like RebaseAtCommit for RebaseAtMerge.
func (r *RepoIO) MergePR(ctx context.Context, payload *defs.RepoIOMergePRPayload) error {
	client, err := Instance().GetClientForInstallationID(payload.InstallationID)
	if err != nil {
		return err
	}

	// Create a copy branch name for the PR which will merge to main
	copy_branch := payload.DefaultBranch + "-copy-for-" + payload.TargetBranch

	// Get the latest commit SHA of the default branch
	commits, _, err := client.Repositories.ListCommits(ctx, payload.RepoOwner, payload.RepoName, &gh.CommitsListOptions{
		SHA: payload.DefaultBranch,
	})
	if err != nil {
		return err
	}

	if len(commits) == 0 {
		return fmt.Errorf("no commits found on branch %s", payload.DefaultBranch)
	}

	// get the latest sha
	sha := *commits[0].SHA

	// Create a new branch based on the latest commit
	ref := &gh.Reference{
		Ref: gh.String("refs/heads/" + copy_branch),
		Object: &gh.GitObject{
			SHA: &sha,
		},
	}

	if _, _, err = client.Git.CreateRef(ctx, payload.RepoOwner, payload.RepoName, ref); err != nil {
		return err
	}

	// Function to perform rebase
	merge := func(base, head, message string) error {
		req := &gh.RepositoryMergeRequest{
			Base:          gh.String(base),
			Head:          gh.String(head),
			CommitMessage: gh.String(message),
		}
		_, _, err = client.Repositories.Merge(ctx, payload.RepoOwner, payload.RepoName, req)

		return err
	}

	// merge target branch with the new branch
	if err = merge(copy_branch, payload.TargetBranch, "Rebasing "+payload.TargetBranch+" with "+copy_branch); err != nil {
		return err
	}

	// merge the new branch with the main branch
	if err = merge(payload.DefaultBranch, copy_branch, "Rebasing "+copy_branch+" with "+payload.DefaultBranch); err != nil {
		return err
	}

	return nil
}
