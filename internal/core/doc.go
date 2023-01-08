package core

import (
	_ "github.com/deepmap/oapi-codegen/pkg/codegen" // Required for code generation
)

//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --config openapi.codegen.yaml openapi.spec.yaml
