package db

import (
	"go.breu.io/quantm/internal/db/config"
	"go.breu.io/quantm/internal/db/entities"
)

// Queries is a wrapper around the config.Queries singleton.
func Queries() *entities.Queries {
	return config.Queries()
}
