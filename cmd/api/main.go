package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"go.breu.io/ctrlplane/internal/conf"
	"go.breu.io/ctrlplane/internal/webhooks"
)

func main() {
	defer conf.Temporal.Client.Close()

	router := chi.NewRouter()

	router.Use(chimiddleware.RequestID)
	router.Use(chimiddleware.RealIP)
	router.Use(chimiddleware.Logger)
	router.Use(chimiddleware.Recoverer)

	router.Get("/webhooks/github", func(response http.ResponseWriter, request *http.Request) {
		fmt.Printf("%+v", request)

		response.Write([]byte("Hello, World!"))
	})

	router.Post("/webhooks/github", webhooks.GithubWebhook)

	http.ListenAndServe(":8000", router)
}

func init() {
	conf.InitService("ctrlplane-api")
	conf.InitKratos()
	conf.InitGithub()
	conf.InitTemporal()
	conf.InitTemporalClient()
}
