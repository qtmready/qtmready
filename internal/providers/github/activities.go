// Copyright © 2023, Breu, Inc. <info@breu.io>. All rights reserved.
//
// This software is made available by Breu, Inc., under the terms of the BREU COMMUNITY LICENSE AGREEMENT, Version 1.0,
// found at https://www.breu.io/license/community. BY INSTALLING, DOWNLOADING, ACCESSING, USING OR DISTRIBUTING ANY OF
// THE SOFTWARE, YOU AGREE TO THE TERMS OF THE LICENSE AGREEMENT.
//
// The above copyright notice and the subsequent license agreement shall be included in all copies or substantial
// portions of the software.
//
// Breu, Inc. HEREBY DISCLAIMS ANY AND ALL WARRANTIES AND CONDITIONS, EXPRESS, IMPLIED, STATUTORY, OR OTHERWISE, AND
// SPECIFICALLY DISCLAIMS ANY WARRANTY OF MERCHANTABILITY OR FITNESS FOR A PARTICULAR PURPOSE, WITH RESPECT TO THE
// SOFTWARE.
//
// Breu, Inc. SHALL NOT BE LIABLE FOR ANY DAMAGES OF ANY KIND, INCLUDING BUT NOT LIMITED TO, LOST PROFITS OR ANY
// CONSEQUENTIAL, SPECIAL, INCIDENTAL, INDIRECT, OR DIRECT DAMAGES, HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY,
// ARISING OUT OF THIS AGREEMENT. THE FOREGOING SHALL APPLY TO THE EXTENT PERMITTED BY APPLICABLE LAW.

package github

import (
	"context"
	"fmt"
	"io"
	"strings"

	gh "github.com/google/go-github/v62/github"
	"go.temporal.io/sdk/activity"

	"go.breu.io/quantm/internal/core"
	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/shared"
)

type (
	// Activities groups all the activities for the github provider.
	Activities struct{}
)

// CreateOrUpdateInstallation creates or update the Installation.
func (a *Activities) CreateOrUpdateInstallation(ctx context.Context, payload *Installation) (*Installation, error) {
	log := activity.GetLogger(ctx)
	installation, err := a.GetInstallation(ctx, payload.InstallationID)

	// if we get the installation, the error will be nil
	if err == nil {
		log.Info("installation found, updating status ...")

		installation.Status = payload.Status
	} else {
		log.Info("installation not found, creating ...")
		log.Debug("payload", "payload", payload)

		installation = payload
	}

	log.Info("saving installation ...")

	if err := db.Save(installation); err != nil {
		log.Error("error saving installation", "error", err)
		return installation, err
	}

	return installation, nil
}

// CreateOrUpdateGithubRepo creates a single row for Repo.
func (a *Activities) CreateOrUpdateGithubRepo(ctx context.Context, payload *Repo) error {
	log := activity.GetLogger(ctx)
	repo, err := a.GetGithubRepo(ctx, payload)

	// if we get the repo, the error will be nil
	if err == nil {
		log.Info("repository found, updating ...")
	} else {
		log.Info("repository not found, creating ...")
		log.Debug("payload", "payload", payload)
	}

	if err := db.Save(repo); err != nil {
		log.Error("error saving repository ...", "error", err)
		return err
	}

	return nil
}

// GetGithubRepo gets Repo against given Repo.
func (a *Activities) GetGithubRepo(ctx context.Context, payload *Repo) (*Repo, error) {
	repo := &Repo{}
	params := db.QueryParams{
		"name":      "'" + payload.Name + "'",
		"full_name": "'" + payload.FullName + "'",
		"github_id": payload.GithubID.String(),
		"team_id":   payload.TeamID.String(),
	}

	if err := db.Get(repo, params); err != nil {
		return payload, err
	}

	return repo, nil
}

// GetInstallation gets Installation against given installation_id.
func (a *Activities) GetInstallation(ctx context.Context, id shared.Int64) (*Installation, error) {
	installation := &Installation{}

	if err := db.Get(installation, db.QueryParams{"installation_id": id.String()}); err != nil {
		return installation, err
	}

	return installation, nil
}

// GetCoreRepo gets entity.Repo against given Repo.
func (a *Activities) GetCoreRepo(ctx context.Context, repo *Repo) (*core.Repo, error) {
	r := &core.Repo{}

	// TODO: add provider name in query
	params := db.QueryParams{
		"provider_id": repo.GithubID.String(),
		"provider":    "'github'",
	}

	if err := db.Get(r, params); err != nil {
		return r, err
	}

	return r, nil
}

// GetCoreRepo gets entity.Stack against given core Repo.
func (a *Activities) GetStack(ctx context.Context, repo *core.Repo) (*core.Stack, error) {
	s := &core.Stack{}

	params := db.QueryParams{
		"id": repo.StackID.String(),
	}

	if err := db.Get(s, params); err != nil {
		return s, err
	}

	return s, nil
}

// GetLatestCommit gets latest commit for default branch of the provided repo.
func (a *Activities) GetLatestCommit(ctx context.Context, payload *core.RepoIOGetLatestCommitPayload) (*core.LatestCommit, error) {
	logger := activity.GetLogger(ctx)
	prepo := &Repo{}

	logger.Info(
		"Starting Activity: GetLatestCommit with ...",
		"repoID", payload.RepoID,
		"branch", payload.BranchName,
		"github_private_key", Instance().PrivateKey,
	)

	if err := db.Get(prepo, db.QueryParams{"github_id": payload.RepoID}); err != nil {
		return nil, err
	}

	client, err := Instance().GetClientFromInstallation(prepo.InstallationID)
	if err != nil {
		logger.Error("GetClientFromInstallation failed", "Error", err)
		return nil, err
	}

	// TODO: move to some genernic function or activity
	repo, _, err := client.Repositories.Get(ctx, strings.Split(prepo.FullName, "/")[0], prepo.Name)
	if err != nil {
		logger.Error("ChangesInBranch Activity", "Error", err)
		return nil, err
	}

	gb, _, err := client.Repositories.
		GetBranch(context.Background(), strings.Split(prepo.FullName, "/")[0], prepo.Name, payload.BranchName, 10)
	if err != nil {
		logger.Error("GetBranch for Github Repo failed", "Error", err)
		return nil, err
	}

	commit := &core.LatestCommit{
		RepoName:  repo.GetName(),
		RepoUrl:   repo.GetHTMLURL(),
		Branch:    *gb.Name,
		SHA:       *gb.Commit.SHA,
		CommitUrl: *gb.Commit.HTMLURL,
	}

	logger.Debug("Repo", "Name", prepo.FullName, "Branch name", gb.Name, "Last commit", commit)

	return commit, nil
}

// TODO - break it to smalller activities (create, delete and merge).
func (a *Activities) RebaseAndMerge(ctx context.Context, payload *core.RepoIORebaseAndMergePayload) (string, error) {
	logger := activity.GetLogger(ctx)

	client, err := Instance().GetClientFromInstallation(payload.InstallationID)
	if err != nil {
		logger.Error("GetClientFromInstallation failed", "Error", err)
		return "", err
	}

	// Get the default branch (e.g., "main")
	// TODO: move to some genernic function or activity
	repo, _, err := client.Repositories.Get(ctx, payload.RepoOwner, payload.RepoName)
	if err != nil {
		logger.Error("RebaseAndMerge Activity", "Error", err)
		return "", err
	}

	defaultBranch := *repo.DefaultBranch
	newBranchName := defaultBranch + "-tempcopy-for-target-" + payload.TargetBranchName

	// Get the latest commit SHA of the default branch
	commits, _, err := client.Repositories.ListCommits(ctx, payload.RepoOwner, payload.RepoName, &gh.CommitsListOptions{
		SHA: defaultBranch,
	})
	if err != nil {
		logger.Error("RebaseAndMerge Activity", "Error", err)
		return "", err
	}

	// // Use the latest commit SHA
	// if len(commits) == 0 {
	// 	shared.Logger().Error("RebaseAndMerge Activity", "No commits found in the default branch.", nil)
	// 	return err.Error(), err
	// }

	latestCommitSHA := *commits[0].SHA

	// Create a new branch based on the latest commit
	ref := &gh.Reference{
		Ref: gh.String("refs/heads/" + newBranchName),
		Object: &gh.GitObject{
			SHA: &latestCommitSHA,
		},
	}

	_, _, err = client.Git.CreateRef(ctx, payload.RepoOwner, payload.RepoName, ref)
	if err != nil {
		logger.Error("RebaseAndMerge Activity", "Error", err)
		return "", err
	}

	logger.Info("RebaseAndMerge Activity", "Branch created successfully: ", newBranchName)

	// Perform rebase of the target branch with the new branch
	rebaseRequest := &gh.RepositoryMergeRequest{
		Base:          &newBranchName,
		Head:          &payload.TargetBranchName,
		CommitMessage: gh.String("Rebasing " + payload.TargetBranchName + " with " + newBranchName),
	}

	_, _, err = client.Repositories.Merge(ctx, payload.RepoOwner, payload.RepoName, rebaseRequest)
	if err != nil {
		logger.Error("RebaseAndMerge Activity", "Error", err)
		return "", err
	}

	logger.Info("RebaseAndMerge Activity", "status",
		fmt.Sprintf("Branch %s rebased with %s successfully.\n", payload.TargetBranchName, newBranchName))

	// Perform rebase of the new branch with the main branch
	rebaseRequest = &gh.RepositoryMergeRequest{
		Base:          &defaultBranch,
		Head:          &newBranchName,
		CommitMessage: gh.String("Rebasing " + newBranchName + " with " + defaultBranch),
	}

	repoCommit, _, err := client.Repositories.Merge(ctx, payload.RepoOwner, payload.RepoName, rebaseRequest)
	if err != nil {
		logger.Error("RebaseAndMerge Activity", "Error", err)
		return err.Error(), err
	}

	logger.Info("RebaseAndMerge Activity", "status",
		fmt.Sprintf("Branch %s rebased with %s successfully.\n", newBranchName, defaultBranch))

	return *repoCommit.SHA, nil
}

func (a *Activities) TriggerCIAction(ctx context.Context, payload *core.RepoIOTriggerCIActionPayload) error {
	logger := activity.GetLogger(ctx)

	logger.Debug("activity TriggerGithubAction started")

	client, err := Instance().GetClientFromInstallation(payload.InstallationID)
	if err != nil {
		logger.Error("GetClientFromInstallation failed", "Error", err)
		return err
	}

	workflowName := "cicd_quantm.yaml" //TODO: either fix this or obtain it somehow

	paylod := gh.CreateWorkflowDispatchEventRequest{
		Ref: payload.TargetBranch,
		Inputs: map[string]any{
			"target-branch": payload.TargetBranch,
		},
	}

	res, err := client.Actions.CreateWorkflowDispatchEventByFileName(ctx, payload.RepoOwner, payload.RepoName, workflowName, paylod)
	if err != nil {
		logger.Error("TriggerGithubAction", "Error", err)
		return err
	}

	logger.Debug("TriggerGithubAction", "response", res)

	return nil
}

func (a *Activities) DeployChangeset(ctx context.Context, payload *core.RepoIODeployChangesetPayload) error {
	logger := activity.GetLogger(ctx)
	logger.Debug("DeployChangeset", "github activity DeployChangeset started for changeset", payload.ChangesetID)

	gh_action_name := "deploy_quantm.yaml" //TODO: fixed it for now

	// get installationID, repoName, repoOwner from github_repos table
	githubRepo := &Repo{}
	params := db.QueryParams{
		"github_id": payload.RepoID,
	}

	if err := db.Get(githubRepo, params); err != nil {
		return err
	}

	client, err := Instance().GetClientFromInstallation(githubRepo.InstallationID)
	if err != nil {
		logger.Error("GetClientFromInstallation failed", "Error", err)
		return err
	}

	paylod := gh.CreateWorkflowDispatchEventRequest{
		Ref: "main",
		Inputs: map[string]any{
			"changesetId": payload.ChangesetID,
		},
	}

	var repoOwner, repoName string

	parts := strings.Split(githubRepo.FullName, "/")

	if len(parts) == 2 {
		repoOwner = parts[0]
		repoName = parts[1]
	}

	res, err := client.Actions.CreateWorkflowDispatchEventByFileName(ctx, repoOwner, repoName, gh_action_name, paylod)
	if err != nil {
		logger.Error("DeployChangeset", "Error", err)
		return err
	}

	logger.Debug("DeployChangeset", "response", res)

	return nil
}

func (a *Activities) TagCommit(ctx context.Context, payload *core.RepoIOTagCommitPayload) error {
	logger := activity.GetLogger(ctx)
	// get installationID, repoName, repoOwner from github_repos table
	githubRepo := &Repo{}
	params := db.QueryParams{
		"github_id": payload.RepoID,
	}

	if err := db.Get(githubRepo, params); err != nil {
		return err
	}

	client, err := Instance().GetClientFromInstallation(githubRepo.InstallationID)
	if err != nil {
		logger.Error("GetClientFromInstallation failed", "Error", err)
		return err
	}

	tag := &gh.Tag{
		Tag:     &payload.TagName,
		Message: &payload.TagMessage,
		Object: &gh.GitObject{
			SHA:  &payload.CommitSHA,
			Type: gh.String("commit"), // Specify the type of object being tagged
		},
	}

	var repoOwner, repoName string

	parts := strings.Split(githubRepo.FullName, "/")

	if len(parts) == 2 {
		repoOwner = parts[0]
		repoName = parts[1]
	}

	_, _, err = client.Git.CreateTag(ctx, repoOwner, repoName, tag)
	if err != nil {
		logger.Error("TagCommit: creating tag", "Error", err)
	}

	// Push the tag to the remote repository
	ref := "refs/tags/" + payload.TagName
	if _, _, err = client.Git.CreateRef(ctx, repoOwner, repoName, &gh.Reference{
		Ref:    &ref,
		Object: &gh.GitObject{SHA: &payload.CommitSHA},
	}); err != nil {
		logger.Error("TagCommit: pushing tag to remote repository", "Error", err)
	}

	return nil
}

func (a *Activities) DeleteBranch(ctx context.Context, payload *core.RepoIODeleteBranchPayload) error {
	logger := activity.GetLogger(ctx)

	// Get github client
	client, err := Instance().GetClientFromInstallation(payload.InstallationID)
	if err != nil {
		logger.Error("GetClientFromInstallation failed", "Error", err)
		return err
	}

	// Create a new branch based on the latest defaultBranch commit
	ref := &gh.Reference{
		Ref: gh.String("refs/heads/" + payload.BranchName),
	}

	// delete the temp ref if its present before
	if _, err := client.Git.DeleteRef(context.Background(), payload.RepoOwner, payload.RepoName, *ref.Ref); err != nil {
		logger.Error("DeleteBranch", "Error deleting ref"+*ref.Ref, err)

		if strings.Contains(err.Error(), "422 Reference does not exist") {
			// if a ref doesnt exist already, dont return the error
			return nil
		}

		return err
	}

	return nil
}

func (a *Activities) CreateBranch(ctx context.Context, payload *core.RepoIOCreateBranchPayload) error {
	logger := activity.GetLogger(ctx)

	// Get github client
	client, err := Instance().GetClientFromInstallation(payload.InstallationID)
	if err != nil {
		logger.Error("GetClientFromInstallation failed", "Error", err)
		return err
	}

	// Create a new branch based on the latest defaultBranch commit
	ref := &gh.Reference{
		Ref: gh.String("refs/heads/" + payload.BranchName),
		Object: &gh.GitObject{
			SHA: &payload.Commit,
		},
	}

	// create new ref
	if _, _, err = client.Git.CreateRef(context.Background(), payload.RepoOwner, payload.RepoName, ref); err != nil {
		logger.Error("CreateBranch activity", "Error", err)
		// dont want to retry this workflow so not returning error, just log and return
		return nil
	}

	return nil
}

func (a *Activities) MergeBranch(ctx context.Context, payload *core.RepoIOMergeBranchPayload) error {
	logger := activity.GetLogger(ctx)

	// Get github client for operations
	client, err := Instance().GetClientFromInstallation(payload.InstallationID)
	if err != nil {
		logger.Error("GetClientFromInstallation failed", "Error", err)
		return err
	}

	// targetBranch will be merged into the baseBranch
	rebaseReq := &gh.RepositoryMergeRequest{
		Base:          &payload.BaseBranch,
		Head:          &payload.TargetBranch,
		CommitMessage: gh.String("Rebasing " + payload.TargetBranch + " with " + payload.BaseBranch),
	}

	if _, _, err := client.Repositories.Merge(context.Background(), payload.RepoOwner, payload.RepoName, rebaseReq); err != nil {
		logger.Error("Merge failed", "Error", err)
		return err
	}

	return nil
}

func (a *Activities) DetectChange(ctx context.Context, payload *core.RepoIODetectChangePayload) (*core.BranchChanges, error) {
	logger := activity.GetLogger(ctx)

	// Get github client for operations
	client, err := Instance().GetClientFromInstallation(payload.InstallationID)
	if err != nil {
		logger.Error("GetClientFromInstallation failed", "Error", err)
		return nil, err
	}

	// TODO: move to some genernic function or activity
	repo, _, err := client.Repositories.Get(ctx, payload.RepoOwner, payload.RepoName)
	if err != nil {
		logger.Error("ChangesInBranch Activity", "Error", err)
		return nil, err
	}

	comparison, _, err := client.Repositories.
		CompareCommits(context.Background(), payload.RepoOwner, payload.RepoName, payload.DefaultBranch, payload.TargetBranch, nil)
	if err != nil {
		logger.Error("Error in ChangesInBranch", "Error", err)
		return nil, err
	}

	var changes, additions, deletions int

	var changedFiles []string

	for _, file := range comparison.Files {
		changes += file.GetChanges()
		additions += file.GetAdditions()
		deletions += file.GetDeletions()
		changedFiles = append(changedFiles, *file.Filename)
	}

	branchChanges := &core.BranchChanges{
		RepoUrl:    repo.GetHTMLURL(),
		Changes:    changes,
		Additions:  additions,
		Deletions:  deletions,
		CompareUrl: comparison.GetHTMLURL(),
		FileCount:  len(changedFiles),
		Files:      changedFiles,
	}

	logger.Debug("ChangesInBranch", "total changes in branch "+payload.TargetBranch, changes)

	return branchChanges, nil
}

func (a *Activities) GetAllBranches(ctx context.Context, payload *core.RepoIOGetAllBranchesPayload) ([]string, error) {
	logger := activity.GetLogger(ctx)

	// get github client
	client, err := Instance().GetClientFromInstallation(payload.InstallationID)
	if err != nil {
		logger.Error("GetClientFromInstallation failed", "Error", err)
		return nil, err
	}

	var branchNames []string

	page := 1

	for {
		branches, resp, err := client.Repositories.ListBranches(ctx, payload.RepoOwner, payload.RepoName, &gh.BranchListOptions{
			ListOptions: gh.ListOptions{
				Page:    page,
				PerPage: 30, // Adjust this value as needed
			},
		})
		if err != nil {
			logger.Error("GetAllBranches: could not get branches", "Error", err)
			return nil, err
		}

		for _, branch := range branches {
			branchNames = append(branchNames, *branch.Name)
		}

		// Check if there are more pages to fetch
		if resp.NextPage == 0 {
			break // No more pages
		}

		page = resp.NextPage
	}

	return branchNames, nil
}

func (a *Activities) GetRepoTeamID(ctx context.Context, payload *core.RepoIOGetRepoTeamIDPayload) (string, error) {
	logger := activity.GetLogger(ctx)
	prepo := &Repo{}

	if err := db.Get(prepo, db.QueryParams{"github_id": payload.RepoID}); err != nil {
		logger.Error("GetRepoTeamID failed", "Error", err)
		return "", err
	}

	logger.Info("GetRepoTeamID Activity", "Get Repo Team ID successfully: ", prepo.TeamID)

	return prepo.TeamID.String(), nil
}

func (a *Activities) GetAllRelevantActions(ctx context.Context, payload *core.RepoIOGetAllRelevantActionsPayload) error {
	logger := activity.GetLogger(ctx)

	// get github client
	client, err := Instance().GetClientFromInstallation(payload.InstallationID)
	if err != nil {
		logger.Error("GetClientFromInstallation failed", "Error", err)
		return err
	}

	// List repository workflows
	workflows, _, err := client.Actions.ListWorkflows(ctx, payload.RepoOwner, payload.RepoName, nil)
	if err != nil {
		return err
	}

	// var labeledWorkflows []string

	// initialize workflow status record map
	actionWorkflowStatuses[payload.RepoName] = make(map[string]string)

	// Iterate through each workflow
	for _, workflow := range workflows.Workflows {
		// Download the content of the workflow file
		content, _, err := client.Repositories.DownloadContents(ctx, payload.RepoOwner, payload.RepoName, *workflow.Path, nil)
		if err != nil {
			return err
		}

		// Read the content bytes
		contentBytes, err := io.ReadAll(content)
		if err != nil {
			return err
		}

		// Convert content bytes to string
		contentStr := string(contentBytes)

		// Check if the workflow is triggered by the specified label
		if strings.Contains(contentStr, "quantm ready") {
			logger.Debug("action file: " + *workflow.Path)

			// labeledWorkflows = append(labeledWorkflows, *workflow.Path)
			actionWorkflowStatuses[payload.RepoName][*workflow.Path] = "idle"
		}
	}

	return nil
}

func (a *Activities) GetRepoByProviderID(
	ctx context.Context, payload *core.RepoIOGetRepoByProviderIDPayload,
) (*core.RepoProviderData, error) {
	prepo := &Repo{}

	// NOTE: these activities are used in api not in temporal workflow use shared.Logger()
	if err := db.Get(prepo, db.QueryParams{"id": payload.ProviderID}); err != nil {
		shared.Logger().Error("GetRepoByProviderID failed", "Error", err)
		return nil, err
	}

	shared.Logger().Info("Get Repo by Provider ID successfully")

	rpd := &core.RepoProviderData{
		Name:          prepo.Name,
		DefaultBranch: prepo.DefaultBranch,
	}

	return rpd, nil
}

func (a *Activities) UpdateRepoHasRarlyWarning(ctx context.Context, payload *core.RepoIOUpdateRepoHasRarlyWarningPayload) error {
	prepo := &Repo{}

	// NOTE: these activities are used in api not in temporal workflow use shared.Logger()
	if err := db.Get(prepo, db.QueryParams{"id": payload.ProviderID}); err != nil {
		shared.Logger().Error("UpdateRepoHasRarlWarning failed", "Error", err)
		return err
	}

	prepo.HasEarlyWarning = true

	if err := db.Save(prepo); err != nil {
		return err
	}

	shared.Logger().Info("Update Repo Has Rarly Warning successfully")

	return nil
}

func (a *Activities) GetOrgUsers(ctx context.Context, payload *core.RepoIOGetOrgUsersPayload) error {
	logger := activity.GetLogger(ctx)
	prepo := &Repo{}

	if err := db.Get(prepo, db.QueryParams{"team_id": payload.TeamID}); err != nil {
		logger.Error("GetOrgUsers: Unable to get the gitub repo", "Error", err)
		return err
	}

	client, err := Instance().GetClientFromInstallation(prepo.InstallationID)
	if err != nil {
		logger.Error("GetClientFromInstallation failed", "Error", err)
		return err
	}

	olopts := &gh.OrganizationsListOptions{
		ListOptions: gh.ListOptions{
			Page:    1,
			PerPage: 1, // Adjust this value as needed
		},
	}

	// NOTE: make dynamic and get the sepecific organizations
	orgs, _, err := client.Organizations.ListAll(ctx, olopts)
	if err != nil {
		logger.Error("List the organizations failed", "Error", err)
		return err
	}

	lmopts := &gh.ListMembersOptions{
		ListOptions: gh.ListOptions{
			Page:    1,
			PerPage: 1, // Adjust this value as needed
		},
	}

	members, _, err := client.Organizations.ListMembers(ctx, *orgs[0].Name, lmopts)
	if err != nil {
		logger.Error("List the organization members failed", "Error", err)
		return err
	}

	logger.Info("organization members count: ", len(members))

	// Save the github org users
	for _, member := range members {
		m := GithubOrgMembers{
			Name:    *member.Name,
			Email:   *member.Email,
			Company: *member.Company,
		}

		if err := db.Save(&m); err != nil {
			logger.Error("Error saving github org members", "Error", err)
			return err
		}
	}

	return nil
}

func (a *Activities) RefreshDefaultBranches(ctx context.Context, payload *core.RepoIORefreshDefaultBranchesPayload) error {
	logger := activity.GetLogger(ctx)

	prepos := make([]*Repo, 0)
	if err := db.Filter(&Repo{}, prepos, db.QueryParams{"team_id": payload.TeamID}); err != nil {
		shared.Logger().Error("Error filter repos", "error", err)
	}

	client, err := Instance().GetClientFromInstallation(prepos[0].InstallationID)
	if err != nil {
		logger.Error("GetClientFromInstallation failed", "Error", err)
		return err
	}

	// Save the github org users
	for _, prepo := range prepos {
		repo, _, err := client.Repositories.Get(ctx, strings.Split(prepo.FullName, "/")[0], prepo.Name)
		if err != nil {
			logger.Error("RefreshDefaultBranches Activity", "Error", err)
			return err
		}

		prepo.DefaultBranch = *repo.DefaultBranch

		if err := db.Save(prepo); err != nil {
			logger.Error("Error saving ggithub repo", "Error", err)
			return err
		}
	}

	return nil
}

// GetRepoForInstallation filters repositories by installation ID and GitHub ID.
// A repo on GitHub can be associated with multiple installations. This function is used to get the repo for a specific installation.
func (a *Activities) GetReposForInstallation(ctx context.Context, installationID, githubID string) ([]Repo, error) {
	var repos []Repo
	err := db.Filter(&Repo{}, &repos, db.QueryParams{
		"installation_id": installationID,
		"github_id":       githubID,
	})

	if err != nil {
		return nil, err
	}

	return repos, nil
}

// GetCoreRepoByProviderID retrieves a core repository by its provider ID.
func (a *Activities) GetCoreRepoByProviderID(ctx context.Context, id string) (*core.Repo, error) {
	repo := &core.Repo{}
	if err := db.Get(repo, db.QueryParams{"provider_id": "'" + id + "'"}); err != nil {
		return nil, err
	}

	return repo, nil
}
