package github

import (
	"io/ioutil"
	"net/http"

	c "go.breu.io/ctrlplane/internal/conf"
	"go.uber.org/zap"
)

// ConsumeGithubInstallationEvent handles GitHub installation events
func webhook(response http.ResponseWriter, request *http.Request) {
	id := request.Header.Get("X-GitHub-Delivery")
	signature := request.Header.Get("X-Hub-Signature")

	if signature == "" {
		handleError(id, ErrorMissingHeaderGithubSignature, http.StatusUnauthorized, response)
		return
	}

	body, _ := ioutil.ReadAll(request.Body)

	if err := Github.VerifyWebhookSignature(body, signature); err != nil {
		handleError(id, err, http.StatusUnauthorized, response)
		return
	}

	headerEvent := request.Header.Get("X-GitHub-Event")

	if headerEvent == "" {
		handleError(id, ErrorMissingHeaderGithubEvent, http.StatusBadRequest, response)
		return
	}

	event := GithubEvent(headerEvent)

	if handle, exists := eventHandlers[event]; exists {
		c.Logger.Info("Received event", zap.String("event", string(event)), zap.String("request_id", id))
		handle(id, body, response)
	} else {
		c.Logger.Error("Unsupported event: " + headerEvent)
		handleError(id, ErrorInvalidEvent, http.StatusBadRequest, response)
	}
}
