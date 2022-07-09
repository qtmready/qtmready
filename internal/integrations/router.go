package integrations

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"go.breu.io/ctrlplane/internal/integrations/github"
)

func Router() http.Handler {
	router := chi.NewRouter()
	router.Mount("github", github.Router())
	return router
}
