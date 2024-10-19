package main

import (
	"net/http"

	"github.com/a-h/rest"

	"go.breu.io/quantm/internal/db/entities"
	"go.breu.io/quantm/internal/shared"
)

func orgs(api *rest.API) {
	api.RegisterModel(
		rest.ModelOf[entities.CreateOrgParams](),
		rest.WithDescription("Create Organization Params"),
	)

	api.RegisterModel(
		rest.ModelOf[entities.UpdateOrgParams](),
		rest.WithDescription("Update Organization Params"),
	)

	api.RegisterModel(
		rest.ModelOf[entities.Org](),
		rest.WithDescription("Organization"),
	)

	// /auth/orgs

	api.Post("/auth/orgs").
		HasDescription("Create Organization").
		HasOperationID("CreateOrg").
		HasRequestModel(rest.ModelOf[entities.CreateOrgParams]()).
		HasResponseModel(http.StatusCreated, rest.ModelOf[entities.Org]()).
		HasResponseModel(http.StatusBadRequest, rest.ModelOf[shared.BadRequest]()).
		HasResponseModel(http.StatusUnauthorized, rest.ModelOf[shared.Unauthorized]()).
		HasResponseModel(http.StatusInternalServerError, rest.ModelOf[shared.InternalServerError]())

	// /auth/orgs/{id}

	api.Put("/auth/orgs/{id}").
		HasDescription("Update Organization").
		HasOperationID("UpdateOrg").
		HasPathParameter("id", rest.PathParam{
			Description: "org id",
			Type:        rest.PrimitiveTypeString,
		}).
		HasRequestModel(rest.ModelOf[entities.UpdateOrgParams]()).
		HasResponseModel(http.StatusOK, rest.ModelOf[entities.Org]()).
		HasResponseModel(http.StatusBadRequest, rest.ModelOf[shared.BadRequest]()).
		HasResponseModel(http.StatusUnauthorized, rest.ModelOf[shared.Unauthorized]()).
		HasResponseModel(http.StatusInternalServerError, rest.ModelOf[shared.InternalServerError]())
}
