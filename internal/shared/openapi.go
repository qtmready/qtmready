package shared

import (
	_ "github.com/deepmap/oapi-codegen/pkg/codegen" // Required for code generation
)

//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen -config openapi.codegen.yaml -package shared -generate types,skip-prune -o types.gen.go openapi.spec.yaml
