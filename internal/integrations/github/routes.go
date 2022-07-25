package github

import (
	"io/ioutil"
	"net/http"

	"go.breu.io/ctrlplane/internal/common"
	"go.breu.io/ctrlplane/internal/common/utils"
	"go.uber.org/zap"
)

// handles the incoming webhook
func webhook(writer http.ResponseWriter, request *http.Request) {
	id := request.Header.Get("X-GitHub-Delivery")
	signature := request.Header.Get("X-Hub-Signature")

	if signature == "" {
		utils.HandleHTTPError(writer, ErrorMissingHeaderGithubSignature, http.StatusUnauthorized)
		return
	}

	body, _ := ioutil.ReadAll(request.Body)

	if err := Github.VerifyWebhookSignature(body, signature); err != nil {
		utils.HandleHTTPError(writer, err, http.StatusUnauthorized)
		return
	}

	headerEvent := request.Header.Get("X-GitHub-Event")

	if headerEvent == "" {
		utils.HandleHTTPError(writer, ErrorMissingHeaderGithubEvent, http.StatusBadRequest)
		return
	}

	event := WebhookEvent(headerEvent)

	// We get the handler for the event. see event_handlers.go
	if handle, exists := eventHandlers[event]; exists {
		common.Logger.Info("Received event", zap.String("event", string(event)), zap.String("request_id", id))
		handle(writer, body, id)
	} else {
		common.Logger.Error("Unsupported event: " + headerEvent)
		utils.HandleHTTPError(writer, ErrorInvalidEvent, http.StatusBadRequest)
	}
}

func completeInstallation(writer http.ResponseWriter, request *http.Request) {}
