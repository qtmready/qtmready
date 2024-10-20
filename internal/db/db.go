package db

import (
	"go.breu.io/quantm/internal/db/config"
	"go.breu.io/quantm/internal/db/entities"
)

func Queries() *entities.Queries {
	return config.Queries()
}
