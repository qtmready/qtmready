package sequel

import (
	"context"

	"github.com/google/uuid"

	"go.breu.io/quantm/internal/db/config"
	"go.breu.io/quantm/internal/db/entities"
	"go.breu.io/quantm/internal/erratic"
	"go.breu.io/quantm/internal/shared"
)

type (
	CreateTeamRequest struct {
		Name string `json:"name" validate:"required"`
	}

	GetTeamRequest struct {
		ID uuid.UUID `json:"id" validate:"required"`
	}

	UpdateTeamRequest = entities.UpdateTeamParams
)

// CreateTeam Create a new team.
func CreateTeam(ctx context.Context, req CreateTeamRequest) (entities.Team, error) {
	if err := shared.Validator().Struct(req); err != nil {
		return entities.Team{}, erratic.NewBadRequestError().FormatValidationError(err)
	}

	val := ctx.Value("org_id")
	if val == nil {
		return entities.Team{}, erratic.NewUnauthorizedError().NotLoggedIn()
	}

	org_id, ok := val.(uuid.UUID)
	if !ok {
		return entities.Team{}, erratic.NewInternalServerError()
	}

	team, err := config.Queries().CreateTeam(ctx, entities.CreateTeamParams{Name: req.Name, OrgID: org_id})
	if err != nil {
		return entities.Team{}, erratic.NewInternalServerError().DataBaseError(err)
	}

	return team, nil
}

// GetTeam Get a team by ID.
func GetTeam(ctx context.Context, req GetTeamRequest) (entities.Team, error) {
	err := shared.Validator().Struct(req)
	if err != nil {
		return entities.Team{}, erratic.NewBadRequestError().FormatValidationError(err)
	}

	val := ctx.Value("org_id")
	if val == nil {
		return entities.Team{}, erratic.NewUnauthorizedError().NotLoggedIn()
	}

	org_id, ok := val.(uuid.UUID)
	if !ok {
		return entities.Team{}, erratic.NewInternalServerError()
	}

	team, err := config.Queries().GetTeam(ctx, req.ID)
	if err != nil {
		return entities.Team{}, erratic.NewNotFoundError("team_id", req.ID.String())
	}

	if team.OrgID != org_id {
		return entities.Team{}, erratic.NewUnauthorizedError().IllegalAccess().AddInfo("team_id", req.ID.String())
	}

	return team, nil
}

// UpdateTeam Update a team by ID.
func UpdateTeam(ctx context.Context, req UpdateTeamRequest) (entities.Team, error) {
	if err := shared.Validator().Struct(req); err != nil {
		return entities.Team{}, erratic.NewBadRequestError().FormatValidationError(err)
	}

	val := ctx.Value("org_id")
	if val == nil {
		return entities.Team{}, erratic.NewUnauthorizedError().NotLoggedIn()
	}

	org_id, ok := val.(uuid.UUID)
	if !ok {
		return entities.Team{}, erratic.NewInternalServerError()
	}

	team, err := config.Queries().UpdateTeam(ctx, req)
	if err != nil {
		return entities.Team{}, erratic.NewInternalServerError().DataBaseError(err)
	}

	if team.OrgID != org_id {
		return entities.Team{}, erratic.NewUnauthorizedError().IllegalAccess().AddInfo("team_id", team.ID.String())
	}

	return team, nil
}
