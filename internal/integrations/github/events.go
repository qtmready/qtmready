package github

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	tc "go.temporal.io/sdk/client"
	"go.uber.org/zap"

	conf "go.breu.io/ctrlplane/internal/common"
)

// A Map of event types to their respective handlers
var eventHandlers = map[GithubEvent]func(string, []byte, http.ResponseWriter){
	GithubInstallationEvent:     handleInstallationEvent,
	GithubAppAuthorizationEvent: handleAuthEvent,
	GithubPushEvent:             handlePushEvent,
}

// handle github installation event
func handleInstallationEvent(id string, body []byte, response http.ResponseWriter) {
	payload := GithubInstallationEventPayload{}
	if err := json.Unmarshal(body, &payload); err != nil {
		handleError(id, ErrorPayloadParser, http.StatusBadRequest, response)
		return
	}

	options := tc.StartWorkflowOptions{
		ID:        id + "::" + strconv.Itoa(int(payload.Installation.ID)),
		TaskQueue: conf.Temporal.Queues.Integrations,
	}

	exe, err := conf.Temporal.Client.ExecuteWorkflow(context.Background(), options, WorkflowOnGithubInstall, payload)

	if err != nil {
		conf.Logger.Error(err.Error())
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
	}

	response.WriteHeader(http.StatusCreated)
	response.Write([]byte(exe.GetRunID()))
}

// handle github app authorization event
func handleAuthEvent(id string, body []byte, response http.ResponseWriter) {
	data, _ := json.MarshalIndent(body, "", "  ")
	conf.Logger.Debug("App authorization event received")
	conf.Logger.Debug(string(data))

	response.WriteHeader(http.StatusCreated)
	response.Write([]byte(""))
}

// handle github push event
func handlePushEvent(id string, body []byte, response http.ResponseWriter) {
	payload := GithubPushEventPayload{}
	if err := json.Unmarshal(body, &payload); err != nil {
		handleError(id, ErrorPayloadParser, http.StatusBadRequest, response)
		return
	}
	data, _ := json.MarshalIndent(payload, "", "  ")
	conf.Logger.Debug("Push event received")
	conf.Logger.Debug(string(data))

	response.WriteHeader(http.StatusCreated)
	response.Write([]byte(""))
}

// handleError handles an error and writes it to the response.
func handleError(id string, err error, status int, response http.ResponseWriter) {
	conf.Logger.Error(err.Error(), zap.String("request_id", id))
	response.WriteHeader(status)
	response.Write([]byte(err.Error()))
}
