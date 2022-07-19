package github

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"go.temporal.io/sdk/client"

	"go.breu.io/ctrlplane/internal/common"
	"go.breu.io/ctrlplane/internal/common/utils"
)

// A Map of event types to their respective handlers
var eventHandlers = map[WebhookEvent]func(string, []byte, http.ResponseWriter){
	InstallationEvent:     handleInstallationEvent,
	AppAuthorizationEvent: handleAuthEvent,
	PushEvent:             handlePushEvent,
}

// handle github installation event
func handleInstallationEvent(id string, body []byte, response http.ResponseWriter) {
	payload := InstallationEventPayload{}
	if err := json.Unmarshal(body, &payload); err != nil {
		utils.HandleHttpError(id, ErrorPayloadParser, http.StatusBadRequest, response)
		return
	}

	opts := client.StartWorkflowOptions{
		ID:        "github.webhooks.installation.id." + strconv.Itoa(int(payload.Installation.ID)) + "." + string(InstallationEvent) + "." + payload.Action,
		TaskQueue: common.Temporal.Queues.Integrations,
	}

	var w *Workflows
	exe, err := common.Temporal.Client.ExecuteWorkflow(context.Background(), opts, w.OnInstall, payload)

	if err != nil {
		utils.HandleHttpError(id, err, http.StatusInternalServerError, response)
	}

	response.WriteHeader(http.StatusCreated)
	response.Write([]byte(exe.GetRunID()))
}

// handle github push event
func handlePushEvent(id string, body []byte, response http.ResponseWriter) {
	payload := PushEventPayload{}
	if err := json.Unmarshal(body, &payload); err != nil {
		utils.HandleHttpError(id, ErrorPayloadParser, http.StatusBadRequest, response)
		return
	}

	opts := client.StartWorkflowOptions{
		ID:        "github.webhooks.integrations.id" + strconv.Itoa(payload.Installation.ID) + "." + string(PushEvent) + ".ref." + payload.Ref,
		TaskQueue: common.Temporal.Queues.Integrations,
	}
	var w *Workflows
	exe, err := common.Temporal.Client.ExecuteWorkflow(context.Background(), opts, w.OnPush, payload)

	if err != nil {
		utils.HandleHttpError(id, err, http.StatusInternalServerError, response)
	}

	response.WriteHeader(http.StatusCreated)
	response.Write([]byte(exe.GetRunID()))
}

// handle github app authorization event
func handleAuthEvent(_ string, body []byte, response http.ResponseWriter) {
	data, _ := json.MarshalIndent(body, "", "  ")
	common.Logger.Debug("App authorization event received")
	common.Logger.Debug(string(data))

	response.WriteHeader(http.StatusCreated)
	response.Write([]byte(""))
}
