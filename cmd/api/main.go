package main

import (
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"go.breu.io/ctrlplane/internal/conf"
	"go.breu.io/ctrlplane/internal/webhooks"
)

var wait sync.WaitGroup

func init() {
	defer wait.Done()
	conf.Service.ReadConf()
	conf.Service.InitLogger()

	conf.EventStream.ReadConf()
	conf.Temporal.ReadConf()
	conf.Github.ReadConf()
	conf.DB.ReadConf()

	wait.Add(3)
	go conf.DB.InitSessionWithRunMigrations()
	go conf.EventStream.InitConnection()
	go conf.Temporal.InitClient()
	wait.Wait()

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

	router.Mount("/webhooks", webhooks.Router())

	http.ListenAndServe(":8000", router)
}
