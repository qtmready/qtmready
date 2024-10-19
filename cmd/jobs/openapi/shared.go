package main

import (
	"github.com/a-h/rest"
	"github.com/getkin/kin-openapi/openapi3"

	"go.breu.io/quantm/internal/shared"
)

type Empty struct{}

func models_shared(api *rest.API) {
	api.RegisterModel(
		rest.ModelOf[Empty](),
		rest.WithDescription("No Content"),
	)

	api.RegisterModel(
		rest.ModelOf[shared.APIError](),
		rest.WithDescription("Default API Error"),
		func(s *openapi3.Schema) {
			code := s.Properties["code"]
			code.Value.Description = "HTTP Status Code"
			code.Value.Type = &openapi3.Types{openapi3.TypeInteger}
			code.Value.WithDefault(500).WithMin(400).WithMax(600)

			message := s.Properties["message"]
			message.Value.Description = "Error Message"
			message.Value.Type = &openapi3.Types{openapi3.TypeString}
			message.Value.WithDefault("Internal Server Error")
		},
	)

	api.RegisterModel(
		rest.ModelOf[shared.BadRequest](),
		rest.WithDescription("Bad Request"),
	)

	api.RegisterModel(
		rest.ModelOf[shared.Unauthorized](),
		rest.WithDescription("Unauthorized"),
	)

	api.RegisterModel(
		rest.ModelOf[shared.NotFound](),
		rest.WithDescription("Not Found"),
	)

	api.RegisterModel(
		rest.ModelOf[shared.InternalServerError](),
		rest.WithDescription("Internal Server Error"),
	)
}
