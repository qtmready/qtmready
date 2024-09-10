//go:build generate
// +build generate

// Package github provides functionality for GitHub provider.
package github

import (
	_ "github.com/deepmap/oapi-codegen/v2/pkg/codegen" // Required for code generation
	_ "gopkg.in/yaml.v2"
)

//go:generate go run github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen --config openapi.codegen.yaml ../../../api/openapi/github/v1/schema.yaml
