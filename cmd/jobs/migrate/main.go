package main

import (
	"sync"

	"go.breu.io/ctrlplane/internal/db"
	"go.breu.io/ctrlplane/internal/shared"
)

func main() {
	waigroup := sync.WaitGroup{}
	// Reading the configuration from the environment
	shared.Service.ReadEnv()
	shared.Service.InitLogger()
	db.DB.ReadEnv()
	// Reading the configuration from the environment ... Done

	shared.Logger.Info("Running Migrations ...", "version", shared.Service.Version())
	waigroup.Add(1)

	go func() {
		defer waigroup.Done()
		db.DB.InitSessionWithMigrations()
	}()

	waigroup.Wait()
	shared.Logger.Info("Migrations Done", "version", shared.Service.Version())
}
