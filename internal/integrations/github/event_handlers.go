package github

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	tc "go.temporal.io/sdk/client"
	"go.uber.org/zap"

	"go.breu.io/ctrlplane/internal/common"
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
		handleError(id, ErrorPayloadParser, http.StatusBadRequest, response)
		return
	}

	options := tc.StartWorkflowOptions{
		ID:        id + "::" + strconv.Itoa(int(payload.Installation.ID)),
		TaskQueue: common.Temporal.Queues.Integrations,
	}

	var w *Workflows
	exe, err := common.Temporal.Client.ExecuteWorkflow(context.Background(), options, w.OnInstallationEvent, payload)

	if err != nil {
		handleError(id, err, http.StatusInternalServerError, response)
	}

	response.WriteHeader(http.StatusCreated)
	response.Write([]byte(exe.GetRunID()))
}

// handle github app authorization event
func handleAuthEvent(id string, body []byte, response http.ResponseWriter) {
	data, _ := json.MarshalIndent(body, "", "  ")
	common.Logger.Debug("App authorization event received")
	common.Logger.Debug(string(data))

	response.WriteHeader(http.StatusCreated)
	response.Write([]byte(""))
}

// handle github push event
func handlePushEvent(id string, body []byte, response http.ResponseWriter) {
	payload := PushEventPayload{}
	if err := json.Unmarshal(body, &payload); err != nil {
		handleError(id, ErrorPayloadParser, http.StatusBadRequest, response)
		return
	}
	data, _ := json.MarshalIndent(payload, "", "  ")
	common.Logger.Debug("Push event received")
	common.Logger.Debug(string(data))

	response.WriteHeader(http.StatusCreated)
	response.Write([]byte(""))
}

// handleError handles an error and writes it to the response.
func handleError(id string, err error, status int, response http.ResponseWriter) {
	common.Logger.Error(err.Error(), zap.String("request_id", id))
	response.WriteHeader(status)
	response.Write([]byte(err.Error()))
}
