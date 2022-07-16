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
		utils.HandleHttpError(id, ErrorMissingHeaderGithubSignature, http.StatusUnauthorized, writer)
		return
	}

	body, _ := ioutil.ReadAll(request.Body)

	if err := Github.VerifyWebhookSignature(body, signature); err != nil {
		utils.HandleHttpError(id, err, http.StatusUnauthorized, writer)
		return
	}

	headerEvent := request.Header.Get("X-GitHub-Event")

	if headerEvent == "" {
		utils.HandleHttpError(id, ErrorMissingHeaderGithubEvent, http.StatusBadRequest, writer)
		return
	}

	event := WebhookEvent(headerEvent)

	if handle, exists := eventHandlers[event]; exists {
		common.Logger.Info("Received event", zap.String("event", string(event)), zap.String("request_id", id))
		handle(id, body, writer)
	} else {
		common.Logger.Error("Unsupported event: " + headerEvent)
		utils.HandleHttpError(id, ErrorInvalidEvent, http.StatusBadRequest, writer)
	}
}

func completeInstallation(writer http.ResponseWriter, request *http.Request) {}
