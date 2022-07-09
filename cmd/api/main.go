package main

import (
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"go.breu.io/ctrlplane/internal/conf"
	"go.breu.io/ctrlplane/internal/integrations"
)

var waiter sync.WaitGroup

func init() {
	conf.Service.ReadConf()
	conf.Service.InitLogger()

	conf.EventStream.ReadConf()
	conf.Temporal.ReadConf()
	integrations.Github.ReadEnv()
	conf.DB.ReadConf()

	waiter.Add(3)

	go func() {
		defer waiter.Done()
		conf.DB.InitSessionWithRunMigrations()
	}()

	go func() {
		defer waiter.Done()
		conf.EventStream.InitConnection()
	}()

	go func() {
		defer waiter.Done()
		conf.Temporal.InitClient()
	}()

	waiter.Wait()

	conf.Logger.Info("Initializing Service ... Done")
}

func main() {
	defer conf.DB.Session.Close()
	defer conf.Temporal.Client.Close()

	router := chi.NewRouter()

	router.Use(chimw.RequestID)
	router.Use(chimw.RealIP)
	router.Use(chimw.Logger)
	router.Use(chimw.Recoverer)

	router.Mount("/integrations", integrations.Router())

	http.ListenAndServe(":8000", router)
}
