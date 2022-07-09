package webhooks

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	tc "go.temporal.io/sdk/client"

	"go.breu.io/ctrlplane/internal/conf"
	"go.breu.io/ctrlplane/internal/temporal/workflows"
	"go.breu.io/ctrlplane/internal/types"
)

// A Map of event types to their respective handlers
var eventHandlers = map[types.GithubEvent]func(string, []byte, http.ResponseWriter){
	types.GithubInstallationEvent:     handleGithubInstallationEvent,
	types.GithubAppAuthorizationEvent: handleGithubAppAuthorizationEvent,
	types.GithubPushEvent:             handleGithubPushEvent,
}

// handle github installation event
func handleGithubInstallationEvent(id string, body []byte, response http.ResponseWriter) {
	payload := &types.GithubInstallationEventPayload{}
	if err := json.Unmarshal(body, payload); err != nil {
		handleError(id, ErrorPayloadParser, http.StatusBadRequest, response)
		return
	}

	options := tc.StartWorkflowOptions{
		ID:        id + "::" + strconv.Itoa(int(payload.Installation.ID)),
		TaskQueue: conf.Temporal.Queues.Webhooks,
	}

	exe, err := conf.Temporal.Client.ExecuteWorkflow(context.Background(), options, workflows.OnGithubInstall, payload)

	if err != nil {
		conf.Logger.Error(err.Error())
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
	}

	response.WriteHeader(http.StatusCreated)
	response.Write([]byte(exe.GetRunID()))
}

// handle github app authorization event
func handleGithubAppAuthorizationEvent(id string, body []byte, response http.ResponseWriter) {
	data, _ := json.MarshalIndent(body, "", "  ")
	conf.Logger.Debug("App authorization event received")
	conf.Logger.Debug(string(data))

	response.WriteHeader(http.StatusCreated)
	response.Write([]byte(""))
}

// handle github push event
func handleGithubPushEvent(id string, body []byte, response http.ResponseWriter) {
	data, _ := json.MarshalIndent(body, "", "  ")
	conf.Logger.Debug("Push event received")
	conf.Logger.Debug(string(data))

	response.WriteHeader(http.StatusCreated)
	response.Write([]byte(""))
}
