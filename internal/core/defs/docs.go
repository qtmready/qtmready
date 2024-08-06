package defs

import (
	_ "github.com/deepmap/oapi-codegen/v2/pkg/codegen" // Required for code generation
	_ "gopkg.in/yaml.v2"
)

// nolint
//go:generate go run github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen --config oapi-codegen.yaml ../../../api/openapi/core/v1/components.yaml
