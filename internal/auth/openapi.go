package auth

import (
	_ "github.com/deepmap/oapi-codegen/pkg/codegen" // Required for code generation
)

//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen -old-config-style -config openapi.codegen.yaml -package auth -generate types -o types.gen.go openapi.spec.yaml
//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen -old-config-style -config openapi.codegen.yaml -package auth -generate server -o server.gen.go openapi.spec.yaml
