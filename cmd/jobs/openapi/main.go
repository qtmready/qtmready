package main

import (
	"encoding/json"
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3gen"
	"go.breu.io/quantm/internal/db/entities"
)

func main() {
	schemas := make(openapi3.Schemas)
	ref, _ := openapi3gen.NewSchemaRefForValue(&entities.CreateOrgParams{}, schemas)

	data, err := json.MarshalIndent(schemas, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Printf("schemas: %s\n", data)

	if data, err = json.MarshalIndent(ref, "", "  "); err != nil {
		panic(err)
	}

	fmt.Printf("schemaRef: %s\n", data)
}
