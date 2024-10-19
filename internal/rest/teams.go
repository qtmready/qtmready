package rest

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/db/entities"
)

type (
	CreateTeamRequest struct {
		Name string `json:"name"`
	}

	GetTeamRequest struct {
		ID uuid.UUID `json:"id"`
	}

	UpdateTeamRequest = entities.UpdateTeamParams
)

// CreateTeam Create a new team
func CreateTeam(ctx context.Context, req CreateTeamRequest) (entities.Team, error) {
	val := ctx.Value("org_id")
	orgID, ok := val.(uuid.UUID)
	if !ok {
		return entities.Team{}, fmt.Errorf("invalid org id")
	}
	return db.Queries().CreateTeam(ctx, entities.CreateTeamParams{Name: req.Name, OrgID: orgID})
}

// GetTeam Get a team by ID
func GetTeam(ctx context.Context, req GetTeamRequest) (entities.Team, error) {
	team, err := db.Queries().GetTeam(ctx, req.ID)
	if err != nil {
		return entities.Team{}, err
	}

	val := ctx.Value("org_id")
	orgID, ok := val.(uuid.UUID)
	if !ok {
		return entities.Team{}, NewUnauthorizedError("team_id", req.ID.String())
	}

	if team.OrgID != orgID {
		return entities.Team{}, NewUnauthorizedError("team_id", req.ID.String())
	}

	return team, nil
}

// UpdateTeam Update a team by ID
func UpdateTeam(ctx context.Context, req UpdateTeamRequest) (entities.Team, error) {
	team, err := db.Queries().UpdateTeam(ctx, req)
	if err != nil {
		return entities.Team{}, err
	}

	val := ctx.Value("org_id")
	orgID, ok := val.(uuid.UUID)
	if !ok {
		return entities.Team{}, NewUnauthorizedError("team_id", req.ID.String())
	}

	if team.OrgID != orgID {
		return entities.Team{}, NewUnauthorizedError("team_id", req.ID.String())
	}

	return team, nil
}
