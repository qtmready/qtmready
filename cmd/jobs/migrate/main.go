package main

import (
	"go.breu.io/quantm/internal/db"
)

func main() {
	db.NewSession(
		db.FromEnvironment(),
		db.WithSessionCreation(),
		db.WithMigrations(),
	)
}
