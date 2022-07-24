package github

import (
	"io/ioutil"
	"net/http"

	"go.breu.io/ctrlplane/internal/common"
	"go.breu.io/ctrlplane/internal/common/utils"
	"go.uber.org/zap"
)

// handles the incoming webhook
func webhook(response http.ResponseWriter, request *http.Request) {
	id := request.Header.Get("X-GitHub-Delivery")
	signature := request.Header.Get("X-Hub-Signature")

	if signature == "" {
		utils.HandleHTTPError(id, ErrorMissingHeaderGithubSignature, http.StatusUnauthorized, response)
		return
	}

	body, _ := ioutil.ReadAll(request.Body)

	if err := Github.VerifyWebhookSignature(body, signature); err != nil {
		utils.HandleHTTPError(id, err, http.StatusUnauthorized, response)
		return
	}

	headerEvent := request.Header.Get("X-GitHub-Event")

	if headerEvent == "" {
		utils.HandleHTTPError(id, ErrorMissingHeaderGithubEvent, http.StatusBadRequest, response)
		return
	}

	event := WebhookEvent(headerEvent)

	// We get the handler for the event. see event_handlers.go
	if handle, exists := eventHandlers[event]; exists {
		common.Logger.Info("Received event", zap.String("event", string(event)), zap.String("request_id", id))
		handle(id, body, response)
	} else {
		common.Logger.Error("Unsupported event: " + headerEvent)
		utils.HandleHTTPError(id, ErrorInvalidEvent, http.StatusBadRequest, response)
	}
}

func completeInstallation(response http.ResponseWriter, request *http.Request) {}
