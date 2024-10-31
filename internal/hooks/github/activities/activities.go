package activities

import (
	"context"

	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/db/entities"
)

type (
	// Entity groups all the activities for the github hook.
	Entity struct{}
)

// --- User and Team and Org Management ---

// GetUserByID retrieves a user from the database by their ID.
func (a *Entity) GetUserByID(ctx context.Context, id db.String) (*entities.User, error) {
	// NOTE - I think not an efficient way. looking for better approach.
	uid, err := id.ToUUID()
	if err != nil {
		return nil, err
	}

	user, err := db.Queries().GetUserByID(ctx, uid)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// SaveUser saves the provided user to the authentication provider.
func (a *Entity) SaveUser(ctx context.Context, arg entities.CreateUserParams) (*entities.User, error) {
	created, err := db.Queries().CreateUser(ctx, arg)
	if err != nil {
		return nil, err
	}

	return &created, nil
}

// SaveTeam saves a new team in the authentication provider.
func (a *Entity) SaveTeam(ctx context.Context, arg entities.CreateTeamParams) (*entities.Team, error) {
	created, err := db.Queries().CreateTeam(ctx, arg)
	if err != nil {
		return nil, err
	}

	return &created, nil
}

// GetTeamByID retrieves a team by its ID.
func (a *Entity) GetTeamByID(ctx context.Context, id db.String) (*entities.Team, error) {
	// NOTE - I think not an efficient way. looking for better approach.
	teamID, err := id.ToUUID()
	if err != nil {
		return nil, err
	}

	team, err := db.Queries().GetTeamByID(ctx, teamID)
	if err != nil {
		return nil, err
	}

	return &team, nil
}

// SaveGithubOrgOrg saves a new github orgnization in the authentication provider.
func (a *Entity) SaveGithubOrg(ctx context.Context, arg entities.CreateGithubOrgParams) (*entities.GithubOrg, error) {
	// TODO - handle no github orgnization use case.
	created, err := db.Queries().CreateGithubOrg(ctx, arg)
	if err != nil {
		return nil, err
	}

	return &created, nil
}

// GetGithubOrgByID etrieves a github org by its ID.
func (a *Entity) GetGithubOrgByID(ctx context.Context, id db.String) (*entities.GithubOrg, error) {
	// NOTE - I think not an efficient way. looking for better approach.
	orgID, err := id.ToUUID()
	if err != nil {
		return nil, err
	}

	org, err := db.Queries().GetGithubOrgByID(ctx, orgID)
	if err != nil {
		return nil, err
	}

	return &org, nil
}

// --- Installation Management ---

// CreateOrUpdateInstallation creates or updates an Installation.
func (a *Entity) CreateOrUpdateInstallation(ctx context.Context, payload *entities.GithubInstallation) error {
	arg := entities.GetGithubInstallationByInstallationIDAndInstallationLoginParams{
		InstallationID:    payload.InstallationID,
		InstallationLogin: payload.InstallationLogin,
	}

	installation, err := db.Queries().GetGithubInstallationByInstallationIDAndInstallationLogin(ctx, arg)
	if err != nil {
		create := entities.CreateGithubInstallationParams{
			InstallationID:      payload.InstallationID,
			InstallationLogin:   payload.InstallationLogin,
			InstallationLoginID: payload.InstallationLoginID,
			InstallationType:    payload.InstallationType,
			SenderID:            payload.SenderID,
			SenderLogin:         payload.SenderLogin,
			Status:              payload.Status,
		}

		if _, err := db.Queries().CreateGithubInstallation(ctx, create); err != nil {
			return err
		}

		return nil
	}

	update := entities.UpdateGithubInstallationParams{
		ID:                  installation.ID,
		InstallationID:      payload.InstallationID,
		InstallationLogin:   payload.InstallationLogin,
		InstallationLoginID: payload.InstallationLoginID,
		InstallationType:    payload.InstallationType,
		SenderID:            payload.SenderID,
		SenderLogin:         payload.SenderLogin,
		Status:              payload.Status,
	}

	if _, err := db.Queries().UpdateGithubInstallation(ctx, update); err != nil {
		return err
	}

	return nil
}

// GetInstallation gets Installation against given installation_id & github login.
func (a *Entity) GetInstallation(
	ctx context.Context, id int64, login string,
) (*entities.GithubInstallation, error) {
	params := entities.GetGithubInstallationByInstallationIDAndInstallationLoginParams{
		InstallationID:    id,
		InstallationLogin: login,
	}

	installation, err := db.Queries().GetGithubInstallationByInstallationIDAndInstallationLogin(ctx, params)

	if err != nil {
		return nil, err
	}

	return &installation, nil
}

// --- Repository Management ---

// CreateOrUpdateGithubRepo creates or updates a single row for Repo.
func (a *Entity) CreateOrUpdateGithubRepo(ctx context.Context, payload *entities.GithubRepo) error {
	params := entities.GetGithubRepoParams{
		Name:     payload.Name,
		FullName: payload.FullName,
		GithubID: payload.GithubID,
	}

	repo, err := db.Queries().GetGithubRepo(ctx, params)
	if err != nil {
		create := entities.CreateGithubRepoParams{
			RepoID:         payload.RepoID,
			InstallationID: payload.InstallationID,
			GithubID:       payload.GithubID,
			Name:           payload.Name,
			FullName:       payload.FullName,
			Url:            payload.Url,
			IsActive:       payload.IsActive,
		}

		if _, err := db.Queries().CreateGithubRepo(ctx, create); err != nil {
			return err
		}

		return nil
	}

	update := entities.UpdateGithubRepoParams{
		ID:             repo.ID,
		RepoID:         payload.RepoID,
		InstallationID: payload.InstallationID,
		GithubID:       payload.GithubID,
		Name:           payload.Name,
		FullName:       payload.FullName,
		Url:            payload.Url,
		IsActive:       payload.IsActive,
	}

	if _, err := db.Queries().UpdateGithubRepo(ctx, update); err != nil {
		return err
	}

	return nil
}

// GetCoreRepo gets entity.Repo against given Repo.
func (*Entity) GetCoreRepo(ctx context.Context, id db.String) (*entities.GetGithubReposWithCoreRepoRow, error) {
	// NOTE - I think not an efficient way. looking for better approach.
	gid, err := id.ToUUID()
	if err != nil {
		return nil, err
	}

	repo, err := db.Queries().GetGithubReposWithCoreRepo(ctx, gid)
	if err != nil {
		return nil, err
	}

	return &repo, nil
}
