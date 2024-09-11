package auth

import (
	"context"

	"go.breu.io/quantm/internal/db"
)

type (
	// TeamUserIO represents the activities for the team user.
	teamuserio struct{}
)

// TeamUserIO creates and returns a new TeamUserIO object.
//
// Example:
//
//	team_user_io := auth.TeamUserIO()
func TeamUserIO() *teamuserio {
	return &teamuserio{}
}

// Get retrieves a team user from the database by their user ID and team ID.
//
// Example:
//
//	team_user, err := auth.TeamUserIO().Get(ctx, user_id, team_id)
func (a *teamuserio) Get(ctx context.Context, user_id, team_id string) (*TeamUser, error) {
	team_user := &TeamUser{}

	return team_user, db.Get(team_user, db.QueryParams{
		"id":      user_id,
		"team_id": team_id,
	})
}

// GetByUserID retrieves a team user from the database by their user ID.
//
// Example:
//
//	team_users, err := auth.TeamUserIO().GetByUserID(ctx, user_id)
func (a *teamuserio) GetByUserID(ctx context.Context, user_id string) ([]TeamUser, error) {
	tus := make([]TeamUser, 0)
	err := db.Filter(&TeamUser{}, tus, db.QueryParams{"id": user_id})

	return tus, err
}

// GetByTeamID retrieves a team user from the database by their team ID.
//
// Example:
//
//	team_users, err := auth.TeamUserIO().GetByTeamID(ctx, team_id)
func (a *teamuserio) GetByTeamID(ctx context.Context, team_id string) ([]TeamUser, error) {
	tus := make([]TeamUser, 0)
	err := db.Filter(&TeamUser{}, tus, db.QueryParams{"team_id": team_id})

	return tus, err
}

// GetByLogin retrieves a team user from the database by their login.
//
// Example:
//
//	team_user, err := auth.TeamUserIO().GetByLogin(ctx, login)
func (a *teamuserio) GetByLogin(ctx context.Context, login string) (*TeamUser, error) {
	team_user := &TeamUser{}

	return team_user, db.Get(team_user, db.QueryParams{"user_login_id": login})
}

// Save creates or updates a team user in the database.
//
// Example:
//
//	team_user, err := auth.TeamUserIO().Save(ctx, team_user)
func (a *teamuserio) Save(ctx context.Context, tu *TeamUser) (*TeamUser, error) {
	fetched := &TeamUser{}

	if err := db.Get(fetched, db.QueryParams{"team_id": tu.TeamID.String(), "user_id": tu.UserID.String()}); err == nil {
		tu.ID = fetched.ID
		tu.CreatedAt = fetched.CreatedAt
	}

	if err := db.Save(tu); err != nil {
		return nil, err
	}

	return tu, nil
}
