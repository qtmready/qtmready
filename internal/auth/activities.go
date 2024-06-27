package auth

import (
	"context"

	"go.breu.io/quantm/internal/db"
)

type (
	Activities struct{}
)

func NewActivities() *Activities {
	return &Activities{}
}

// GetUser retrieves a user from the database based on the provided query parameters.
//
// The function takes a context.Context and db.QueryParams as input, and returns a pointer to a User
// struct and an error. If an error occurs during the database retrieval, it is returned.
func (a *Activities) GetUser(ctx context.Context, params db.QueryParams) (*User, error) {
	user := &User{}
	if err := db.Get(user, params); err != nil {
		return nil, err
	}

	return user, nil
}

// SaveUser saves the provided user to the database.
// It returns the saved user or an error if the save operation failed.
func (a *Activities) SaveUser(ctx context.Context, user *User) (*User, error) {
	if err := db.Save(user); err != nil {
		return nil, err
	}

	return user, nil
}

// CreateTeam creates a new team in the database.
//
// The function takes a context.Context and a pointer to a Team struct as input.
// It attempts to create the team in the database using the db.Create() function.
// If the creation is successful, the function returns the created Team struct and a nil error.
// If there is an error creating the team, the function returns a nil Team struct and the error.
func (a *Activities) CreateTeam(ctx context.Context, team *Team) (*Team, error) {
	if err := db.Create(team); err != nil {
		return nil, err
	}

	return team, nil
}

// GetTeam retrieves a Team from the database using the provided query parameters.
//
// The returned Team pointer should not be modified directly, as it is a reference to the database object.
// If an error occurs during the database query, it will be returned.
func (a *Activities) GetTeam(ctx context.Context, params db.QueryParams) (*Team, error) {
	team := &Team{}
	if err := db.Get(team, params); err != nil {
		return nil, err
	}

	return team, nil
}

// CreateOrUpdateTeamUser creates or updates a team user in the database.
// It takes a TeamUser payload and returns the updated TeamUser.
// If the TeamUser already exists, it updates the IsAdmin, IsActive, and Role fields.
// If the TeamUser does not exist, it creates a new one.
// The function returns the updated or created TeamUser, and any error that occurred.
func (a *Activities) CreateOrUpdateTeamUser(ctx context.Context, payload *TeamUser) (*TeamUser, error) {
	temp := &TeamUser{}

	if err := db.Get(temp, db.QueryParams{"team_id": payload.TeamID.String(), "user_id": payload.UserID.String()}); err == nil {
		payload.ID = temp.ID
		payload.CreatedAt = temp.CreatedAt
	}

	if err := db.Save(payload); err != nil {
		return nil, err
	}

	return payload, nil
}

// Get team user by login.
func (a *Activities) GetTeamUser(ctx context.Context, loginID string) (*TeamUser, error) {
	temp := &TeamUser{}

	if err := db.Get(temp, db.QueryParams{"user_login_id": loginID}); err != nil {
		return nil, err
	}

	return temp, nil
}
