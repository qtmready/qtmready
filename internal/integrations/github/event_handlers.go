package github

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"go.temporal.io/sdk/client"

	"go.breu.io/ctrlplane/internal/cmn"
	"go.breu.io/ctrlplane/internal/cmn/utils"
)

type eventHandler func(writer http.ResponseWriter, payload []byte, id string)

// A Map of event types to their respective handlers
var eventHandlers = map[WebhookEvent]eventHandler{
	InstallationEvent:     handleInstallationEvent,
	AppAuthorizationEvent: handleAuthEvent,
	PushEvent:             handlePushEvent,
}

// handles GitHub installation event
func handleInstallationEvent(writer http.ResponseWriter, body []byte, id string) {
	payload := InstallationEventPayload{}
	if err := json.Unmarshal(body, &payload); err != nil {
		utils.HandleHTTPError(writer, ErrorPayloadParser, http.StatusBadRequest)
		return
	}

	opts := client.StartWorkflowOptions{
		ID:        "github.webhooks.installation.id." + strconv.Itoa(int(payload.Installation.ID)) + "." + string(InstallationEvent) + "." + payload.Action,
		TaskQueue: cmn.Temporal.Queues.Integrations,
	}

	var w *Workflows
	exe, err := cmn.Temporal.Client.ExecuteWorkflow(context.Background(), opts, w.OnInstall, payload)

	if err != nil {
		utils.HandleHTTPError(writer, err, http.StatusInternalServerError)
	}

	writer.WriteHeader(http.StatusCreated)
	writer.Write([]byte(exe.GetRunID()))
}

// handles GitHub push event
func handlePushEvent(writer http.ResponseWriter, body []byte, id string) {
	payload := PushEventPayload{}
	if err := json.Unmarshal(body, &payload); err != nil {
		utils.HandleHTTPError(writer, ErrorPayloadParser, http.StatusBadRequest)
		return
	}

	opts := client.StartWorkflowOptions{
		ID:        "github.webhooks.integrations.id." + strconv.Itoa(payload.Installation.ID) + "." + string(PushEvent) + ".ref." + payload.Ref,
		TaskQueue: cmn.Temporal.Queues.Integrations,
	}
	var w *Workflows
	exe, err := cmn.Temporal.Client.ExecuteWorkflow(context.Background(), opts, w.OnPush, payload)

	if err != nil {
		utils.HandleHTTPError(writer, err, http.StatusInternalServerError)
	}

	writer.WriteHeader(http.StatusCreated)
	writer.Write([]byte(exe.GetRunID()))
}

// handles github app authorization event
func handleAuthEvent(writer http.ResponseWriter, body []byte, id string) {
	data, _ := json.MarshalIndent(body, "", "  ")
	cmn.Log.Debug("App authorization event received")
	cmn.Log.Debug(string(data))

	writer.WriteHeader(http.StatusCreated)
	writer.Write([]byte(""))
}
