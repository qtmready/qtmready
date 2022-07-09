package main

import (
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"go.breu.io/ctrlplane/internal/common"
	"go.breu.io/ctrlplane/internal/db"
	"go.breu.io/ctrlplane/internal/integrations"
	"go.breu.io/ctrlplane/internal/integrations/github"
)

var waiter sync.WaitGroup

func init() {
	common.Service.ReadConf()
	common.Service.InitLogger()

	common.EventStream.ReadConf()
	common.Temporal.ReadEnv()
	github.Github.ReadEnv()
	db.DB.ReadEnv()

	waiter.Add(3)

	go func() {
		defer waiter.Done()
		db.DB.InitSessionWithMigrations()
	}()

	go func() {
		defer waiter.Done()
		common.EventStream.InitConnection()
	}()

	go func() {
		defer waiter.Done()
		common.Temporal.InitClient()
	}()

	waiter.Wait()

	common.Logger.Info("Initializing Service ... Done")
}

func main() {
	defer db.DB.Session.Close()
	defer common.Temporal.Client.Close()

	router := chi.NewRouter()

	router.Use(chimw.RequestID)
	router.Use(chimw.RealIP)
	router.Use(chimw.Logger)
	router.Use(chimw.Recoverer)

	router.Mount("/integrations", integrations.Router())

	http.ListenAndServe(":8000", router)
}
