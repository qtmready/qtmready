package sequel

import (
	"context"

	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/db/entities"
	"go.breu.io/quantm/internal/erratic"
	"go.breu.io/quantm/internal/shared"
)

type (
	CreateOrgRequest struct {
		Name string `json:"name" validate:"required"`
	}

	UpdateOrgRequest = entities.UpdateOrgParams
)

// CreateOrg Create a new organization.
func CreateOrg(ctx context.Context, req CreateOrgRequest) (entities.Org, error) {
	if err := shared.Validator().Struct(req); err != nil {
		return entities.Org{}, erratic.NewBadRequestError().FormatValidationError(err)
	}

	org, err := db.Queries().CreateOrg(ctx, entities.CreateOrgParams{Name: req.Name, Slug: db.CreateSlug(req.Name)})
	if err != nil {
		return entities.Org{}, erratic.NewInternalServerError().DataBaseError(err)
	}

	return org, nil
}

// UpdateOrg Update an existing organization.
func UpdateOrg(ctx context.Context, req UpdateOrgRequest) (entities.Org, error) {
	org, err := db.Queries().UpdateOrg(ctx, req)
	if err != nil {
		return entities.Org{}, erratic.NewInternalServerError().DataBaseError(err)
	}

	return org, nil
}
