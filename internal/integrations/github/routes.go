package github

import (
	"io"
	"net/http"

	"go.breu.io/ctrlplane/internal/cmn/utils"
)

// handles the incoming webhook
func webhook(writer http.ResponseWriter, request *http.Request) {
	id := request.Header.Get("X-GitHub-Delivery")
	signature := request.Header.Get("X-Hub-Signature")

	if signature == "" {
		utils.HandleHTTPError(writer, ErrorMissingHeaderGithubSignature, http.StatusUnauthorized)
		return
	}

	body, _ := io.ReadAll(request.Body)

	if err := Github.VerifyWebhookSignature(body, signature); err != nil {
		http.Error(writer, err.Error(), http.StatusUnauthorized)
		return
	}

	headerEvent := request.Header.Get("X-GitHub-Event")

	if headerEvent == "" {
		http.Error(writer, ErrorMissingHeaderGithubEvent.Error(), http.StatusBadRequest)
		return
	}

	event := WebhookEvent(headerEvent)

	// We get the handler for the event. see event_handlers.go
	if handle, exists := eventHandlers[event]; exists {
		handle(writer, body, id)
	} else {
		http.Error(writer, ErrorInvalidEvent.Error(), http.StatusBadRequest)
		return
	}
}

func completeInstallation(writer http.ResponseWriter, request *http.Request) {}
