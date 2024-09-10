//go:build generate
// +build generate

// Package slack provides functionality for slack provider.
package slack

import (
	_ "github.com/deepmap/oapi-codegen/v2/pkg/codegen" // Required for code generation
	_ "gopkg.in/yaml.v2"
)

//go:generate go run github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen --config openapi.codegen.yaml ../../../api/openapi/slack/v1/schema.yaml
