//go:build generate
// +build generate

// auth providies the authentication and authorization.
package auth

import (
	_ "github.com/deepmap/oapi-codegen/v2/pkg/codegen" // Required for code generation
	_ "gopkg.in/yaml.v2"
)

//go:generate go run github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen --config openapi.codegen.yaml ../../api/openapi/auth/v1/schema.yaml
