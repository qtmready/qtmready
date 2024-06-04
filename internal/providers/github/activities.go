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

	"go.temporal.io/sdk/activity"

	"go.breu.io/quantm/internal/auth"
	"go.breu.io/quantm/internal/core"
	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/shared"
)

type (
	// Activities groups all the activities for the github provider.
	Activities struct{}
)

var (
	authacts *auth.Activities
)

// GetUserByID retrieves a user from the database by their ID.
// The context.Context parameter is used for cancellation and timeouts.
// The id parameter is the unique identifier of the user to retrieve.
// Returns the retrieved user and any error that occurred.
func (a *Activities) GetUserByID(ctx context.Context, id string) (*auth.User, error) {
	params := db.QueryParams{"id": id}

	return authacts.GetUser(ctx, params)
}

// SaveUser saves the provided user to the authentication provider.
func (a *Activities) SaveUser(ctx context.Context, user *auth.User) (*auth.User, error) {
	return authacts.SaveUser(ctx, user)
}

// CreateTeam creates a new team in the authentication provider.
func (a *Activities) CreateTeam(ctx context.Context, team *auth.Team) (*auth.Team, error) {
	return authacts.CreateTeam(ctx, team)
}

// GetTeamByID retrieves a team by its ID.
// ctx is the context for the operation.
// id is the ID of the team to retrieve.
// Returns the retrieved team, or an error if the team could not be found or retrieved.
func (a *Activities) GetTeamByID(ctx context.Context, id string) (*auth.Team, error) {
	params := db.QueryParams{"id": id}

	return authacts.GetTeam(ctx, params)
}

// CreateMemberships creates a new team membership for the given user and team.
// If the user is already a member of the team, the membership is updated to reflect the provided admin status.
// If the user is not already a member of the organization associated with the team, a new organization membership is created.
func (a *Activities) CreateMemberships(ctx context.Context, payload *CreateMembershipsPayload) error {
	orgusr := &OrgUser{}
	teamuser := &auth.TeamUser{
		TeamID:   payload.UserID,
		UserID:   payload.TeamID,
		IsActive: true,
		IsAdmin:  payload.IsAdmin,
	}

	if _, err := authacts.CreateOrUpdateTeamUser(ctx, teamuser); err != nil {
		return err
	}

	params := db.QueryParams{
		"user_id":        payload.UserID.String(),
		"github_user_id": payload.GithubUserID.String(),
		"github_org_id":  payload.GithubOrgID.String(),
	}

	if err := db.Get(orgusr, params); err != nil {
		orgusr.GithubOrgID = payload.GithubOrgID
		orgusr.GithubUserID = payload.GithubUserID
		orgusr.GithubOrgName = payload.GithubOrgName
		orgusr.UserID = payload.UserID

		if err := db.Save(orgusr); err != nil {
			return err
		}
	}

	return nil
}

// CreateOrUpdateInstallation creates or update the Installation.
func (a *Activities) CreateOrUpdateInstallation(ctx context.Context, payload *Installation) (*Installation, error) {
	installation, err := a.GetInstallation(ctx, payload.InstallationID, payload.InstallationLogin)

	// if we get the installation, the error will be nil
	if err == nil {
		installation.Status = payload.Status
	} else {
		installation = payload
	}

	if err := db.Save(installation); err != nil {
		return installation, err
	}

	return installation, nil
}

// GetInstallation gets Installation against given installation_id & github login.
func (a *Activities) GetInstallation(ctx context.Context, id shared.Int64, login string) (*Installation, error) {
	installation := &Installation{}
	params := db.QueryParams{"installation_id": id.String(), "installation_login": login}

	if err := db.Get(installation, params); err != nil {
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

		payload.ID = repo.ID
		payload.CreatedAt = repo.CreatedAt

		repo = payload
	} else {
		log.Info("repository not found, creating ...")
		log.Debug("payload", "payload", payload)
	}

	if err := db.Save(repo); err != nil {
		log.Error("error saving repository ...", "error", err)
		return err
	}

	log.Info("repository saved successfully ...")

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

func (a *Activities) GetRepoByProviderID(
	ctx context.Context, payload *core.RepoIOGetRepoByProviderIDPayload,
) (*core.RepoProviderData, error) {
	repo := &Repo{}

	// NOTE: these activities are used in api not in temporal workflow use shared.Logger()
	if err := db.Get(repo, db.QueryParams{"id": payload.ProviderID}); err != nil {
		shared.Logger().Error("GetRepoByProviderID failed", "Error", err)
		return nil, err
	}

	shared.Logger().Info("Get Repo by Provider ID successfully")

	data := &core.RepoProviderData{
		Name:          repo.Name,
		DefaultBranch: repo.DefaultBranch,
	}

	return data, nil
}

func (a *Activities) UpdateRepoHasRarlyWarning(ctx context.Context, payload *core.RepoIOUpdateRepoHasRarlyWarningPayload) error {
	repo := &Repo{}

	if err := db.Get(repo, db.QueryParams{"id": payload.ProviderID}); err != nil {
		shared.Logger().Error("UpdateRepoHasRarlWarning failed", "Error", err)
		return err
	}

	repo.HasEarlyWarning = true

	if err := db.Save(repo); err != nil {
		return err
	}

	shared.Logger().Info("Update Repo Has Rarly Warning successfully")

	return nil
}

// SyncReposFromGithub syncs repos from github.
// TODO: We will get rate limiting errors here because of when we scale.
// TODO: if the repo has has_early_warning, we will need to update core repo too.
func (a *Activities) SyncReposFromGithub(ctx context.Context, payload *SyncReposFromGithubPayload) error {
	repos := make([]Repo, 0)
	params := db.QueryParams{"installation_id": payload.InstallationID.String(), "team_id": payload.TeamID.String()}

	if err := db.Filter(&Repo{}, &repos, params); err != nil {
		return err
	}

	if client, err := Instance().GetClientForInstallationID(payload.InstallationID); err != nil {
		return err
	} else {
		for idx := range repos {
			repo := repos[idx]

			result, _, err := client.Repositories.Get(ctx, payload.Owner, repo.Name) // TODO: We use use ListReposByOrg here!
			if err != nil {
				return err
			}

			repo.DefaultBranch = result.GetDefaultBranch()

			if err := db.Save(&repo); err != nil {
				return err
			}
		}
	}

	return nil
}

// func (a *Activities) RefreshDefaultBranches(ctx context.Context, payload *core.RepoIORefreshDefaultBranchesPayload) error {
// 	logger := activity.GetLogger(ctx)

// 	repos := make([]Repo, 0)
// 	if err := db.Filter(&Repo{}, &repos, db.QueryParams{"team_id": payload.TeamID}); err != nil {
// 		shared.Logger().Error("Error filter repos", "error", err)
// 		return err
// 	}

// 	logger.Info("provider repos length", "info", len(repos))

// 	client, err := Instance().GetClientForInstallationID(repos[0].InstallationID)
// 	if err != nil {
// 		logger.Error("GetClientFromInstallation failed", "error", err)
// 		return err
// 	}

// 	// Save the github org users
// 	for idx := range repos {
// 		repo := repos[idx]

// 		result, _, err := client.Repositories.Get(ctx, strings.Split(repo.FullName, "/")[0], repo.Name)
// 		if err != nil {
// 			logger.Error("RefreshDefaultBranches Activity", "error", err)
// 			return err
// 		}

// 		repo.DefaultBranch = result.GetDefaultBranch()

// 		if err := db.Save(&repo); err != nil {
// 			logger.Error("Error saving github repo", "error", err)
// 			return err
// 		}
// 	}

// 	return nil
// }

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

// GetCoreRepoByCtrlID retrieves a core repository given the db id of the github repository.
func (a *Activities) GetCoreRepoByCtrlID(ctx context.Context, id string) (*core.Repo, error) {
	repo := &core.Repo{}
	if err := db.Get(repo, db.QueryParams{"ctrl_id": id}); err != nil {
		return nil, err
	}

	return repo, nil
}

// SignalCoreRepoCtrl signals the core repository control workflow with the given signal and payload.
func (a *Activities) SignalCoreRepoCtrl(ctx context.Context, repo *core.Repo, signal shared.WorkflowSignal, payload any) error {
	opts := shared.Temporal().
		Queue(shared.CoreQueue).
		WorkflowOptions(
			shared.WithWorkflowBlock("repo"),
			shared.WithWorkflowBlockID(repo.ID.String()),
		)

	rw := &core.RepoWorkflows{}

	_, err := shared.Temporal().
		Client().
		SignalWithStartWorkflow(context.Background(), opts.ID, signal.String(), payload, opts, rw.RepoCtrl, repo)

	return err
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

// // GetLatestCommit gets latest commit for default branch of the provided repo.
// func (a *Activities) GetLatestCommit(ctx context.Context, payload *core.RepoIOGetLatestCommitPayload) (*core.LatestCommit, error) {
// 	logger := activity.GetLogger(ctx)
// 	prepo := &Repo{}

// 	logger.Info(
// 		"Starting Activity: GetLatestCommit with ...",
// 		"repoID", payload.RepoID,
// 		"branch", payload.BranchName,
// 	)

// 	if err := db.Get(prepo, db.QueryParams{"github_id": payload.RepoID}); err != nil {
// 		return nil, err
// 	}

// 	client, err := Instance().GetClientFromInstallation(prepo.InstallationID)
// 	if err != nil {
// 		logger.Error("GetClientFromInstallation failed", "Error", err)
// 		return nil, err
// 	}

// 	// TODO: move to some genernic function or activity
// 	repo, _, err := client.Repositories.Get(ctx, strings.Split(prepo.FullName, "/")[0], prepo.Name)
// 	if err != nil {
// 		logger.Error("ChangesInBranch Activity", "Error", err)
// 		return nil, err
// 	}

// 	gb, _, err := client.Repositories.
// 		GetBranch(context.Background(), strings.Split(prepo.FullName, "/")[0], prepo.Name, payload.BranchName, 10)
// 	if err != nil {
// 		logger.Error("GetBranch for Github Repo failed", "Error", err)
// 		return nil, err
// 	}

// 	commit := &core.LatestCommit{
// 		RepoName:  repo.GetName(),
// 		RepoUrl:   repo.GetHTMLURL(),
// 		Branch:    *gb.Name,
// 		SHA:       *gb.Commit.SHA,
// 		CommitUrl: *gb.Commit.HTMLURL,
// 	}

// 	logger.Debug("Repo", "Name", prepo.FullName, "Branch name", gb.Name, "Last commit", commit)

// 	return commit, nil
// }

// // TODO - break it to smalller activities (create, delete and merge).
// func (a *Activities) RebaseAndMerge(ctx context.Context, payload *core.RepoIORebaseAndMergePayload) (string, error) {
// 	logger := activity.GetLogger(ctx)

// 	client, err := Instance().GetClientFromInstallation(payload.InstallationID)
// 	if err != nil {
// 		logger.Error("GetClientFromInstallation failed", "Error", err)
// 		return "", err
// 	}

// 	// Get the default branch (e.g., "main")
// 	// TODO: move to some genernic function or activity
// 	repo, _, err := client.Repositories.Get(ctx, payload.RepoOwner, payload.RepoName)
// 	if err != nil {
// 		logger.Error("RebaseAndMerge Activity", "Error", err)
// 		return "", err
// 	}

// 	defaultBranch := *repo.DefaultBranch
// 	newBranchName := defaultBranch + "-tempcopy-for-target-" + payload.TargetBranchName

// 	// Get the latest commit SHA of the default branch
// 	commits, _, err := client.Repositories.ListCommits(ctx, payload.RepoOwner, payload.RepoName, &gh.CommitsListOptions{
// 		SHA: defaultBranch,
// 	})
// 	if err != nil {
// 		logger.Error("RebaseAndMerge Activity", "Error", err)
// 		return "", err
// 	}

// 	// // Use the latest commit SHA
// 	// if len(commits) == 0 {
// 	// 	shared.Logger().Error("RebaseAndMerge Activity", "No commits found in the default branch.", nil)
// 	// 	return err.Error(), err
// 	// }

// 	latestCommitSHA := *commits[0].SHA

// 	// Create a new branch based on the latest commit
// 	ref := &gh.Reference{
// 		Ref: gh.String("refs/heads/" + newBranchName),
// 		Object: &gh.GitObject{
// 			SHA: &latestCommitSHA,
// 		},
// 	}

// 	_, _, err = client.Git.CreateRef(ctx, payload.RepoOwner, payload.RepoName, ref)
// 	if err != nil {
// 		logger.Error("RebaseAndMerge Activity", "Error", err)
// 		return "", err
// 	}

// 	logger.Info("RebaseAndMerge Activity", "Branch created successfully: ", newBranchName)

// 	// Perform rebase of the target branch with the new branch
// 	rebaseRequest := &gh.RepositoryMergeRequest{
// 		Base:          &newBranchName,
// 		Head:          &payload.TargetBranchName,
// 		CommitMessage: gh.String("Rebasing " + payload.TargetBranchName + " with " + newBranchName),
// 	}

// 	_, _, err = client.Repositories.Merge(ctx, payload.RepoOwner, payload.RepoName, rebaseRequest)
// 	if err != nil {
// 		logger.Error("RebaseAndMerge Activity", "Error", err)
// 		return "", err
// 	}

// 	logger.Info("RebaseAndMerge Activity", "status",
// 		fmt.Sprintf("Branch %s rebased with %s successfully.\n", payload.TargetBranchName, newBranchName))

// 	// Perform rebase of the new branch with the main branch
// 	rebaseRequest = &gh.RepositoryMergeRequest{
// 		Base:          &defaultBranch,
// 		Head:          &newBranchName,
// 		CommitMessage: gh.String("Rebasing " + newBranchName + " with " + defaultBranch),
// 	}

// 	repoCommit, _, err := client.Repositories.Merge(ctx, payload.RepoOwner, payload.RepoName, rebaseRequest)
// 	if err != nil {
// 		logger.Error("RebaseAndMerge Activity", "Error", err)
// 		return err.Error(), err
// 	}

// 	logger.Info("RebaseAndMerge Activity", "status",
// 		fmt.Sprintf("Branch %s rebased with %s successfully.\n", newBranchName, defaultBranch))

// 	return *repoCommit.SHA, nil
// }

// func (a *Activities) TriggerCIAction(ctx context.Context, payload *core.RepoIOTriggerCIActionPayload) error {
// 	logger := activity.GetLogger(ctx)

// 	logger.Debug("activity TriggerGithubAction started")

// 	client, err := Instance().GetClientFromInstallation(payload.InstallationID)
// 	if err != nil {
// 		logger.Error("GetClientFromInstallation failed", "Error", err)
// 		return err
// 	}

// 	workflowName := "cicd_quantm.yaml" //TODO: either fix this or obtain it somehow

// 	paylod := gh.CreateWorkflowDispatchEventRequest{
// 		Ref: payload.TargetBranch,
// 		Inputs: map[string]any{
// 			"target-branch": payload.TargetBranch,
// 		},
// 	}

// 	res, err := client.Actions.CreateWorkflowDispatchEventByFileName(ctx, payload.RepoOwner, payload.RepoName, workflowName, paylod)
// 	if err != nil {
// 		logger.Error("TriggerGithubAction", "Error", err)
// 		return err
// 	}

// 	logger.Debug("TriggerGithubAction", "response", res)

// 	return nil
// }

// func (a *Activities) DeployChangeset(ctx context.Context, payload *core.RepoIODeployChangesetPayload) error {
// 	logger := activity.GetLogger(ctx)
// 	logger.Debug("DeployChangeset", "github activity DeployChangeset started for changeset", payload.ChangesetID)

// 	gh_action_name := "deploy_quantm.yaml" //TODO: fixed it for now

// 	// get installationID, repoName, repoOwner from github_repos table
// 	githubRepo := &Repo{}
// 	params := db.QueryParams{
// 		"github_id": payload.RepoID,
// 	}

// 	if err := db.Get(githubRepo, params); err != nil {
// 		return err
// 	}

// 	client, err := Instance().GetClientFromInstallation(githubRepo.InstallationID)
// 	if err != nil {
// 		logger.Error("GetClientFromInstallation failed", "Error", err)
// 		return err
// 	}

// 	paylod := gh.CreateWorkflowDispatchEventRequest{
// 		Ref: "main",
// 		Inputs: map[string]any{
// 			"changesetId": payload.ChangesetID,
// 		},
// 	}

// 	var repoOwner, repoName string

// 	parts := strings.Split(githubRepo.FullName, "/")

// 	if len(parts) == 2 {
// 		repoOwner = parts[0]
// 		repoName = parts[1]
// 	}

// 	res, err := client.Actions.CreateWorkflowDispatchEventByFileName(ctx, repoOwner, repoName, gh_action_name, paylod)
// 	if err != nil {
// 		logger.Error("DeployChangeset", "Error", err)
// 		return err
// 	}

// 	logger.Debug("DeployChangeset", "response", res)

// 	return nil
// }

// func (a *Activities) TagCommit(ctx context.Context, payload *core.RepoIOTagCommitPayload) error {
// 	logger := activity.GetLogger(ctx)
// 	// get installationID, repoName, repoOwner from github_repos table
// 	githubRepo := &Repo{}
// 	params := db.QueryParams{
// 		"github_id": payload.RepoID,
// 	}

// 	if err := db.Get(githubRepo, params); err != nil {
// 		return err
// 	}

// 	client, err := Instance().GetClientFromInstallation(githubRepo.InstallationID)
// 	if err != nil {
// 		logger.Error("GetClientFromInstallation failed", "Error", err)
// 		return err
// 	}

// 	tag := &gh.Tag{
// 		Tag:     &payload.TagName,
// 		Message: &payload.TagMessage,
// 		Object: &gh.GitObject{
// 			SHA:  &payload.CommitSHA,
// 			Type: gh.String("commit"), // Specify the type of object being tagged
// 		},
// 	}

// 	var repoOwner, repoName string

// 	parts := strings.Split(githubRepo.FullName, "/")

// 	if len(parts) == 2 {
// 		repoOwner = parts[0]
// 		repoName = parts[1]
// 	}

// 	_, _, err = client.Git.CreateTag(ctx, repoOwner, repoName, tag)
// 	if err != nil {
// 		logger.Error("TagCommit: creating tag", "Error", err)
// 	}

// 	// Push the tag to the remote repository
// 	ref := "refs/tags/" + payload.TagName
// 	if _, _, err = client.Git.CreateRef(ctx, repoOwner, repoName, &gh.Reference{
// 		Ref:    &ref,
// 		Object: &gh.GitObject{SHA: &payload.CommitSHA},
// 	}); err != nil {
// 		logger.Error("TagCommit: pushing tag to remote repository", "Error", err)
// 	}

// 	return nil
// }

// func (a *Activities) DeleteBranch(ctx context.Context, payload *core.RepoIODeleteBranchPayload) error {
// 	logger := activity.GetLogger(ctx)

// 	// Get github client
// 	client, err := Instance().GetClientFromInstallation(payload.InstallationID)
// 	if err != nil {
// 		logger.Error("GetClientFromInstallation failed", "Error", err)
// 		return err
// 	}

// 	// Create a new branch based on the latest defaultBranch commit
// 	ref := &gh.Reference{
// 		Ref: gh.String("refs/heads/" + payload.BranchName),
// 	}

// 	// delete the temp ref if its present before
// 	if _, err := client.Git.DeleteRef(context.Background(), payload.RepoOwner, payload.RepoName, *ref.Ref); err != nil {
// 		logger.Error("DeleteBranch", "Error deleting ref"+*ref.Ref, err)

// 		if strings.Contains(err.Error(), "422 Reference does not exist") {
// 			// if a ref doesnt exist already, dont return the error
// 			return nil
// 		}

// 		return err
// 	}

// 	return nil
// }

// func (a *Activities) CreateBranch(ctx context.Context, payload *core.RepoIOCreateBranchPayload) error {
// 	logger := activity.GetLogger(ctx)

// 	// Get github client
// 	client, err := Instance().GetClientFromInstallation(payload.InstallationID)
// 	if err != nil {
// 		logger.Error("GetClientFromInstallation failed", "Error", err)
// 		return err
// 	}

// 	// Create a new branch based on the latest defaultBranch commit
// 	ref := &gh.Reference{
// 		Ref: gh.String("refs/heads/" + payload.BranchName),
// 		Object: &gh.GitObject{
// 			SHA: &payload.Commit,
// 		},
// 	}

// 	// create new ref
// 	if _, _, err = client.Git.CreateRef(context.Background(), payload.RepoOwner, payload.RepoName, ref); err != nil {
// 		logger.Error("CreateBranch activity", "Error", err)
// 		// dont want to retry this workflow so not returning error, just log and return
// 		return nil
// 	}

// 	return nil
// }

// func (a *Activities) MergeBranch(ctx context.Context, payload *core.RepoIOMergeBranchPayload) error {
// 	logger := activity.GetLogger(ctx)

// 	// Get github client for operations
// 	client, err := Instance().GetClientFromInstallation(payload.InstallationID)
// 	if err != nil {
// 		logger.Error("GetClientFromInstallation failed", "Error", err)
// 		return err
// 	}

// 	// targetBranch will be merged into the baseBranch
// 	rebaseReq := &gh.RepositoryMergeRequest{
// 		Base:          &payload.BaseBranch,
// 		Head:          &payload.TargetBranch,
// 		CommitMessage: gh.String("Rebasing " + payload.TargetBranch + " with " + payload.BaseBranch),
// 	}

// 	if _, _, err := client.Repositories.Merge(context.Background(), payload.RepoOwner, payload.RepoName, rebaseReq); err != nil {
// 		logger.Error("Merge failed", "Error", err)
// 		return err
// 	}

// 	return nil
// }

// func (a *Activities) DetectChange(ctx context.Context, payload *core.RepoIODetectChangePayload) (*core.BranchChanges, error) {
// 	logger := activity.GetLogger(ctx)

// 	// Get github client for operations
// 	client, err := Instance().GetClientFromInstallation(payload.InstallationID)
// 	if err != nil {
// 		logger.Error("GetClientFromInstallation failed", "Error", err)
// 		return nil, err
// 	}

// 	// TODO: move to some genernic function or activity
// 	repo, _, err := client.Repositories.Get(ctx, payload.RepoOwner, payload.RepoName)
// 	if err != nil {
// 		logger.Error("ChangesInBranch Activity", "Error", err)
// 		return nil, err
// 	}

// 	comparison, _, err := client.Repositories.
// 		CompareCommits(context.Background(), payload.RepoOwner, payload.RepoName, payload.DefaultBranch, payload.TargetBranch, nil)
// 	if err != nil {
// 		logger.Error("Error in ChangesInBranch", "Error", err)
// 		return nil, err
// 	}

// 	var changes, additions, deletions int

// 	var changedFiles []string

// 	for _, file := range comparison.Files {
// 		changes += file.GetChanges()
// 		additions += file.GetAdditions()
// 		deletions += file.GetDeletions()
// 		changedFiles = append(changedFiles, *file.Filename)
// 	}

// 	branchChanges := &core.BranchChanges{
// 		RepoUrl:    repo.GetHTMLURL(),
// 		Changes:    changes,
// 		Additions:  additions,
// 		Deletions:  deletions,
// 		CompareUrl: comparison.GetHTMLURL(),
// 		FileCount:  len(changedFiles),
// 		Files:      changedFiles,
// 	}

// 	logger.Debug("ChangesInBranch", "total changes in branch "+payload.TargetBranch, changes)

// 	return branchChanges, nil
// }

// func (a *Activities) GetAllBranches(ctx context.Context, payload *core.RepoIOGetAllBranchesPayload) ([]string, error) {
// 	logger := activity.GetLogger(ctx)

// 	// get github client
// 	client, err := Instance().GetClientFromInstallation(payload.InstallationID)
// 	if err != nil {
// 		logger.Error("GetClientFromInstallation failed", "Error", err)
// 		return nil, err
// 	}

// 	var branchNames []string

// 	page := 1

// 	for {
// 		branches, resp, err := client.Repositories.ListBranches(ctx, payload.RepoOwner, payload.RepoName, &gh.BranchListOptions{
// 			ListOptions: gh.ListOptions{
// 				Page:    page,
// 				PerPage: 30, // Adjust this value as needed
// 			},
// 		})
// 		if err != nil {
// 			logger.Error("GetAllBranches: could not get branches", "Error", err)
// 			return nil, err
// 		}

// 		for _, branch := range branches {
// 			branchNames = append(branchNames, *branch.Name)
// 		}

// 		// Check if there are more pages to fetch
// 		if resp.NextPage == 0 {
// 			break // No more pages
// 		}

// 		page = resp.NextPage
// 	}

// 	return branchNames, nil
// }

// func (a *Activities) GetRepoTeamID(ctx context.Context, payload *core.RepoIOGetRepoTeamIDPayload) (string, error) {
// 	logger := activity.GetLogger(ctx)
// 	prepo := &Repo{}

// 	if err := db.Get(prepo, db.QueryParams{"github_id": payload.RepoID}); err != nil {
// 		logger.Error("GetRepoTeamID failed", "Error", err)
// 		return "", err
// 	}

// 	logger.Info("GetRepoTeamID Activity", "Get Repo Team ID successfully: ", prepo.TeamID)

// 	return prepo.TeamID.String(), nil
// }

// func (a *Activities) GetAllRelevantActions(ctx context.Context, payload *core.RepoIOGetAllRelevantActionsPayload) error {
// 	logger := activity.GetLogger(ctx)

// 	// get github client
// 	client, err := Instance().GetClientFromInstallation(payload.InstallationID)
// 	if err != nil {
// 		logger.Error("GetClientFromInstallation failed", "Error", err)
// 		return err
// 	}

// 	// List repository workflows
// 	workflows, _, err := client.Actions.ListWorkflows(ctx, payload.RepoOwner, payload.RepoName, nil)
// 	if err != nil {
// 		return err
// 	}

// 	// var labeledWorkflows []string

// 	// initialize workflow status record map
// 	actionWorkflowStatuses[payload.RepoName] = make(map[string]string)

// 	// Iterate through each workflow
// 	for _, workflow := range workflows.Workflows {
// 		// Download the content of the workflow file
// 		content, _, err := client.Repositories.DownloadContents(ctx, payload.RepoOwner, payload.RepoName, *workflow.Path, nil)
// 		if err != nil {
// 			return err
// 		}

// 		// Read the content bytes
// 		contentBytes, err := io.ReadAll(content)
// 		if err != nil {
// 			return err
// 		}

// 		// Convert content bytes to string
// 		contentStr := string(contentBytes)

// 		// Check if the workflow is triggered by the specified label
// 		if strings.Contains(contentStr, "quantm ready") {
// 			logger.Debug("action file: " + *workflow.Path)

// 			// labeledWorkflows = append(labeledWorkflows, *workflow.Path)
// 			actionWorkflowStatuses[payload.RepoName][*workflow.Path] = "idle"
// 		}
// 	}

// 	return nil
// }
