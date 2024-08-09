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
	"strings"

	gh "github.com/google/go-github/v62/github"
	"go.temporal.io/sdk/activity"

	"go.breu.io/quantm/internal/auth"
	"go.breu.io/quantm/internal/core/code"
	"go.breu.io/quantm/internal/core/defs"
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

// GetTeamByID retrieves a user by github user id.
// ctx is the context for the operation.
// id is the ID of the login id provided by github.
// Returns the retrieved team user with message provider data if user not exist return the error and team_user both nil.
func (a *Activities) GetTeamUserByLoginID(ctx context.Context, loginID string) (*auth.TeamUser, error) {
	teamuser, err := authacts.GetTeamUser(ctx, loginID)

	if err != nil {
		// Check if the error message is "not found" and handle accordingly
		// return error nil if the user not found in the system (not connect with message povider)
		// in not foun case we return both err and teamuser to run the workflow but not send message to user
		if strings.Contains(err.Error(), "not found") {
			return nil, nil
		}

		return nil, err
	}

	return teamuser, nil
}

// CreateMemberships creates a new team membership for the given user and team.
// If the user is already a member of the team, the membership is updated to reflect the provided admin status.
// If the user is not already a member of the organization associated with the team, a new organization membership is created.
func (a *Activities) CreateMemberships(ctx context.Context, payload *CreateMembershipsPayload) error {
	orgusr := &OrgUser{}
	teamuser := &auth.TeamUser{
		TeamID:                  payload.TeamID,
		UserID:                  payload.UserID,
		IsActive:                true,
		IsMessageProviderLinked: false,
		IsAdmin:                 payload.IsAdmin,
		UserLoginId:             payload.GithubUserID,
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
func (a *Activities) GetCoreRepo(ctx context.Context, repo *Repo) (*defs.Repo, error) {
	r := &defs.Repo{}

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
	ctx context.Context, payload *defs.RepoIOGetRepoByProviderIDPayload,
) (*defs.RepoProviderData, error) {
	repo := &Repo{}

	// NOTE: these activities are used in api not in temporal workflow use shared.Logger()
	if err := db.Get(repo, db.QueryParams{"id": payload.ProviderID}); err != nil {
		shared.Logger().Error("GetRepoByProviderID failed", "Error", err)
		return nil, err
	}

	shared.Logger().Info("Get Repo by Provider ID successfully")

	data := &defs.RepoProviderData{
		Name:          repo.Name,
		DefaultBranch: repo.DefaultBranch,
	}

	return data, nil
}

func (a *Activities) UpdateRepoHasRarlyWarning(ctx context.Context, payload *defs.RepoIOGetRepoByProviderIDPayload) error {
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

// SyncOrgUsersFromGithub syncs orgainzation users from github.
// NOTE - working only for public org members
// TODO - ifor private org mambers.
func (a *Activities) SyncOrgUsersFromGithub(ctx context.Context, payload *SyncOrgUsersFromGithubPayload) error {
	if client, err := Instance().GetClientForInstallationID(payload.InstallationID); err != nil {
		return err
	} else {
		lmopts := &gh.ListMembersOptions{
			ListOptions: gh.ListOptions{},
		}

		members, _, err := client.Organizations.ListMembers(ctx, payload.GithubOrgName, lmopts)
		if err != nil {
			return err
		}

		for _, member := range members {
			orgusr := &OrgUser{}
			filter := db.QueryParams{
				"github_org_id":  payload.GithubOrgID.String(),
				"github_user_id": shared.Int64(*member.ID).String(),
			}

			// TODO - need to refine
			if err := db.Get(orgusr, filter); err != nil {
				shared.Logger().Debug("member => err", "debug", err)
			}

			orgusr.GithubOrgName = payload.GithubOrgName
			orgusr.GithubOrgID = payload.GithubOrgID
			orgusr.GithubUserID = shared.Int64(*member.ID)

			if err := db.Save(orgusr); err != nil {
				return err
			}
		}
	}

	return nil
}

// func (a *Activities) RefreshDefaultBranches(ctx context.Context, payload *defs.RepoIORefreshDefaultBranchesPayload) error {
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
func (a *Activities) GetCoreRepoByCtrlID(ctx context.Context, id string) (*defs.Repo, error) {
	repo := &defs.Repo{}
	if err := db.Get(repo, db.QueryParams{"ctrl_id": id}); err != nil {
		return nil, err
	}

	return repo, nil
}

// SignalCoreRepoCtrl signals the core repository control workflow with the given signal and payload.
func (a *Activities) SignalCoreRepoCtrl(ctx context.Context, repo *defs.Repo, signal shared.WorkflowSignal, payload any) error {
	opts := shared.Temporal().
		Queue(shared.CoreQueue).
		WorkflowOptions(
			shared.WithWorkflowBlock("repo"),
			shared.WithWorkflowBlockID(repo.ID.String()),
		)

	_, err := shared.Temporal().
		Client().
		SignalWithStartWorkflow(context.Background(), opts.ID, signal.String(), payload, opts, code.RepoCtrl, repo)

	return err
}

func (a *Activities) GithubWorkflowInfo(ctx context.Context, payload *defs.RepoIOWorkflowActionPayload) (*defs.RepoIOWorkflowInfo, error) {
	client, err := Instance().GetClientForInstallationID(payload.InstallationID)
	if err != nil {
		return nil, err
	}

	// List repository workflows
	workflows, _, err := client.Actions.ListWorkflows(ctx, payload.RepoOwner, payload.RepoName, nil)
	if err != nil {
		return nil, err
	}

	// Initialize the result struct
	winfo := &defs.RepoIOWorkflowInfo{
		TotalCount: shared.Int64(workflows.GetTotalCount()),
		Workflows:  make([]*defs.RepIOWorkflow, 0, workflows.GetTotalCount()),
	}

	// Iterate through each workflow
	for _, workflow := range workflows.Workflows {
		detail := &defs.RepIOWorkflow{
			ID:      shared.Int64(*workflow.ID),
			NodeID:  workflow.GetNodeID(),
			Name:    workflow.GetName(),
			Path:    workflow.GetPath(),
			State:   workflow.GetState(),
			HTMLURL: workflow.GetHTMLURL(),
		}

		// Add the workflow to the slice
		winfo.Workflows = append(winfo.Workflows, detail)
	}

	return winfo, nil
}
