package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"go.breu.io/ctrlplane/internal/conf"
	"go.breu.io/ctrlplane/internal/webhooks"
)

func init() {
	conf.ReadSvcConfig("web::api")
	conf.ReadKratosConfig()
	conf.ReadGithubConfig()
	conf.ReadDBConfig()
	conf.InitDBSession()
	conf.ReadTemporalConfig()
	conf.InitTemporalClient()

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

	router.Post("/webhooks/github", webhooks.GithubWebhook)

	http.ListenAndServe(":8000", router)
}
