// Crafted with ❤ at Breu, Inc. <info@breu.io>, Copyright © 2022, 2024.
//
// Functional Source License, Version 1.1, Apache 2.0 Future License
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
	"log/slog"

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

// GetUserByID retrieves a user from the database by their ID.
// The context.Context parameter is used for cancellation and timeouts.
// The id parameter is the unique identifier of the user to retrieve.
// Returns the retrieved user and any error that occurred.
func (a *Activities) GetUserByID(ctx context.Context, id string) (*auth.User, error) {
	params := db.QueryParams{"id": id}

	return auth.UserIO().Get(ctx, params)
}

// SaveUser saves the provided user to the authentication provider.
func (a *Activities) SaveUser(ctx context.Context, user *auth.User) (*auth.User, error) {
	return auth.UserIO().Save(ctx, user)
}

// CreateTeam creates a new team in the authentication provider.
func (a *Activities) CreateTeam(ctx context.Context, team *auth.Team) (*auth.Team, error) {
	return auth.TeamIO().Save(ctx, team)
}

// GetTeamByID retrieves a team by its ID.
// ctx is the context for the operation.
// id is the ID of the team to retrieve.
// Returns the retrieved team, or an error if the team could not be found or retrieved.
func (a *Activities) GetTeamByID(ctx context.Context, id string) (*auth.Team, error) {
	return auth.TeamIO().GetByID(ctx, id)
}

// GetTeamByID retrieves a user by github user id.
// ctx is the context for the operation.
// id is the ID of the login id provided by github.
// Returns the retrieved team user with message provider data if user not exist return the error and team_user both nil.
func (a *Activities) GetTeamUserByLoginID(ctx context.Context, loginID string) (*auth.TeamUser, error) {
	teamuser, err := auth.TeamUserIO().GetByLogin(ctx, loginID)

	if err != nil {
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

	if _, err := auth.TeamUserIO().Save(ctx, teamuser); err != nil {
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
func (a *Activities) GetInstallation(ctx context.Context, id db.Int64, login string) (*Installation, error) {
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

	// NOTE: these activities are used in api not in temporal workflow use slog
	if err := db.Get(repo, db.QueryParams{"id": payload.ProviderID}); err != nil {
		slog.Error("GetRepoByProviderID failed", "Error", err)
		return nil, err
	}

	slog.Info("Get Repo by Provider ID successfully")

	data := &defs.RepoProviderData{
		Name:          repo.Name,
		DefaultBranch: repo.DefaultBranch,
	}

	return data, nil
}

func (a *Activities) UpdateRepoHasRarlyWarning(ctx context.Context, payload *defs.RepoIOGetRepoByProviderIDPayload) error {
	repo := &Repo{}

	if err := db.Get(repo, db.QueryParams{"id": payload.ProviderID}); err != nil {
		slog.Error("UpdateRepoHasRarlWarning failed", "Error", err)
		return err
	}

	repo.HasEarlyWarning = true

	if err := db.Save(repo); err != nil {
		return err
	}

	slog.Info("Update Repo Has Rarly Warning successfully")

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
				"github_user_id": db.Int64(*member.ID).String(),
			}

			// TODO - need to refine
			if err := db.Get(orgusr, filter); err != nil {
				slog.Debug("member => err", "debug", err)
			}

			orgusr.GithubOrgName = payload.GithubOrgName
			orgusr.GithubOrgID = payload.GithubOrgID
			orgusr.GithubUserID = db.Int64(*member.ID)

			if err := db.Save(orgusr); err != nil {
				return err
			}
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

// GetCoreRepoByCtrlID retrieves a core repository given the db id of the github repository.
func (a *Activities) GetCoreRepoByCtrlID(ctx context.Context, id string) (*defs.Repo, error) {
	repo := &defs.Repo{}
	if err := db.Get(repo, db.QueryParams{"ctrl_id": id}); err != nil {
		return nil, err
	}

	return repo, nil
}

// SignalCoreRepoCtrl signals the core repository control workflow with the given signal and payload.
func (a *Activities) SignalCoreRepoCtrl(ctx context.Context, repo *defs.Repo, signal defs.Signal, payload any) error {
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
		TotalCount: db.Int64(workflows.GetTotalCount()),
		Workflows:  make([]*defs.RepIOWorkflow, 0, workflows.GetTotalCount()),
	}

	// Iterate through each workflow
	for _, workflow := range workflows.Workflows {
		w := &defs.RepIOWorkflow{
			ID:      db.Int64(*workflow.ID),
			NodeID:  workflow.GetNodeID(),
			Name:    workflow.GetName(),
			Path:    workflow.GetPath(),
			State:   workflow.GetState(),
			HTMLURL: workflow.GetHTMLURL(),
		}

		winfo.Workflows = append(winfo.Workflows, w)
	}

	return winfo, nil
}
