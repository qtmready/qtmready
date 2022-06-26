package main

import (
	"fmt"
	"net/http"

	_chi "github.com/go-chi/chi/v5"
	_chiMiddleware "github.com/go-chi/chi/v5/middleware"

	"go.breu.io/ctrlplane/cmd/api/middlewares"
	"go.breu.io/ctrlplane/internal/defaults"
	"go.breu.io/ctrlplane/internal/webhooks"
)

func main() {
	defaults.Logger.Info("Starting API")

	router := _chi.NewRouter()

	router.Use(_chiMiddleware.RequestID)
	router.Use(_chiMiddleware.RealIP)
	router.Use(_chiMiddleware.Logger)
	router.Use(middlewares.KratosMiddleware)
	router.Use(_chiMiddleware.Recoverer)

	router.Get("/webhooks/github", func(response http.ResponseWriter, request *http.Request) {
		fmt.Printf("%+v", request)

		response.Write([]byte("Hello, World!"))
	})

	router.Post("/webhooks/github", webhooks.GithubWebhook)

	http.ListenAndServe(":8000", router)
}
