//go:build generate
// +build generate

// shared contains shared code between the various services.
package shared

import (
	_ "github.com/deepmap/oapi-codegen/v2/pkg/codegen" // Required for code generation
	_ "gopkg.in/yaml.v2"
)

//go:generate go run github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen -config openapi.codegen.yaml -package shared -generate types,skip-prune,client -o types.gen.go ../../api/openapi/shared/v1/schema.yaml
