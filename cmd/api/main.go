package main

import (
	"fmt"
	"log"
	"net/http"

	_chi "github.com/go-chi/chi/v5"
	_chiMiddleware "github.com/go-chi/chi/v5/middleware"

	"go.breu.io/ctrlplane/internal/common"
)

func main() {
	logger := log.New(log.Writer(), "", 0)
	logger.Printf("%+v", common.Conf.Github)

	router := _chi.NewRouter()

	router.Use(_chiMiddleware.RequestID)
	router.Use(_chiMiddleware.RealIP)
	router.Use(_chiMiddleware.Logger)
	router.Use(_chiMiddleware.Recoverer)

	// router.Get("/", func(response http.ResponseWriter, request *http.Request) {
	// 	response.Write([]byte("Hello, World!"))
	// })

	router.Post("/webhooks/github", func(response http.ResponseWriter, request *http.Request) {
		event := request.Header.Get("X-GitHub-Event")
		fmt.Println(event)
		response.Write([]byte("Hello, World!"))
	})

	http.ListenAndServe(":3000", router)
}
