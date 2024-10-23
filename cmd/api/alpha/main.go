package main

import (
	"net/http"

	"go.breu.io/quantm/internal/nomad/handler"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle(handler.NewHealthCheckServiceHandler())

	// http.ListenAndServe(":8000", mux)
}
