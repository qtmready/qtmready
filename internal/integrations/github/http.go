package github

import (
	"io/ioutil"
	"net/http"

	"go.breu.io/ctrlplane/internal/common"
	"go.uber.org/zap"
)

// handles the incoming webhook
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
		common.Logger.Info("Received event", zap.String("event", string(event)), zap.String("request_id", id))
		handle(id, body, response)
	} else {
		common.Logger.Error("Unsupported event: " + headerEvent)
		handleError(id, ErrorInvalidEvent, http.StatusBadRequest, response)
	}
}

func completeInstallation(response http.ResponseWriter, request *http.Request) {}
