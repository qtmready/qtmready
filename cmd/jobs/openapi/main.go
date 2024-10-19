package main

import (
	"os"

	"github.com/a-h/rest"
	"gopkg.in/yaml.v3"
)

func main() {
	api := rest.NewAPI("Quantm REST API", rest.WithApplyCustomSchemaToType(custom))
	api.StripPkgPaths = []string{
		"go.breu.io/quantm/internal/db/entities",
		"go.breu.io/quantm/internal/shared",
		"main",
	}

	models_shared(api)

	orgs(api)
	// teams(api)

	spec, err := api.Spec()
	if err != nil {
		panic(err)
	}

	spec.Info.Description = "Quantm REST API"
	spec.Info.Version = "0.1.0"

	enc := yaml.NewEncoder(os.Stdout)
	enc.SetIndent(2)

	if err := enc.Encode(spec); err != nil {
		panic(err)
	}
}
