package webhooks

import (
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"go.breu.io/ctrlplane/internal/conf"
	"go.breu.io/ctrlplane/internal/types"
)

func Router() http.Handler {
	router := chi.NewRouter()
	router.Post("/github", Github)
	return router
}

// ConsumeGithubInstallationEvent handles GitHub installation events
func Github(response http.ResponseWriter, request *http.Request) {
	id := request.Header.Get("X-GitHub-Delivery")
	signature := request.Header.Get("X-Hub-Signature")

	if signature == "" {
		handleError(id, ErrorMissingHeaderGithubSignature, http.StatusUnauthorized, response)
		return
	}

	body, _ := ioutil.ReadAll(request.Body)

	if err := verifySignature(body, signature); err != nil {
		handleError(id, err, http.StatusUnauthorized, response)
		return
	}

	headerEvent := request.Header.Get("X-GitHub-Event")

	if headerEvent == "" {
		handleError(id, ErrorMissingHeaderGithubEvent, http.StatusBadRequest, response)
		return
	}

	event := types.GithubEvent(headerEvent)

	if handle, exists := eventHandlers[event]; exists {
		conf.Logger.Info("Received event", zap.String("event", string(event)), zap.String("request_id", id))
		handle(id, body, response)
	} else {
		conf.Logger.Error("Unsupported event: " + headerEvent)
		handleError(id, ErrorInvalidEvent, http.StatusBadRequest, response)
	}
}
