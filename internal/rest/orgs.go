package rest

import (
	"context"

	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/db/entities"
)

type (
	CreateOrgRequest struct {
		Name string `json:"name"`
	}

	UpdateOrgRequest = entities.UpdateOrgParams
)

// CreateOrg Create a new organization.
func CreateOrg(ctx context.Context, req CreateOrgRequest) (entities.Org, error) {
	return db.Queries().CreateOrg(ctx, entities.CreateOrgParams{Name: req.Name, Slug: db.CreateSlug(req.Name)})
}

// UpdateOrg Update an existing organization.
func UpdateOrg(ctx context.Context, req UpdateOrgRequest) (entities.Org, error) {
	return db.Queries().UpdateOrg(ctx, req)
}
