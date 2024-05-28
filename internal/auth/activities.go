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
