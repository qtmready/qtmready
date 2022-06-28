package webhooks

import (
	"context"
	"encoding/json"
	"net/http"

	_sdkClient "go.temporal.io/sdk/client"

	"go.breu.io/ctrlplane/internal/defaults"
	"go.breu.io/ctrlplane/internal/models"
	"go.breu.io/ctrlplane/internal/workflows"
)

func consumeGithubInstallationEvent(payload models.GithubInstallationEventPayload, response http.ResponseWriter) {
	data, _ := json.Marshal(payload)
	options := _sdkClient.StartWorkflowOptions{
		ID:        string(rune(payload.Installation.ID)),
		TaskQueue: defaults.Conf.Temporal.QUEUES.Webhooks,
	}
	defaults.Logger.Debug("Installation event received")
	defaults.Logger.Debug(string(data))

	exe, err := defaults.Conf.Temporal.Client.ExecuteWorkflow(context.Background(), options, workflows.OnGithubInstall, payload)

	if err != nil {
		defaults.Logger.Error(err.Error())
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
	}

	response.WriteHeader(http.StatusCreated)
	response.Write([]byte(exe.GetRunID()))
}
