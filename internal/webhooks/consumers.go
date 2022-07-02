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

func consumeGithubInstallationEvent(payload types.GithubInstallationEventPayload, response http.ResponseWriter) {
	options := tc.StartWorkflowOptions{
		ID:        strconv.Itoa(int(payload.Installation.ID)),
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

func consumeGithubAppAuthorizationEvent(payload types.GithubAppAuthorizationEventPayload, response http.ResponseWriter) {
	data, _ := json.MarshalIndent(payload, "", "  ")
	conf.Logger.Debug("App authorization event received")
	conf.Logger.Debug(string(data))

	response.WriteHeader(http.StatusCreated)
	response.Write([]byte(""))
}
