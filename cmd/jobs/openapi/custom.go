package main

import (
	"reflect"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/google/uuid"
)

func custom(t reflect.Type, s *openapi3.Schema) {
	if t == reflect.TypeOf(uuid.UUID{}) {
		s.Type = &openapi3.Types{openapi3.TypeString}
		s.Format = "uuid"
		s.Nullable = false
		s.Items = nil
	}
}
