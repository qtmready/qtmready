package webhooks

import (
	"encoding/json"
	"net/http"

	"go.breu.io/ctrlplane/internal/defaults"
)

func consumeGithubInstallationEvent(payload GithubInstallationEventPayload, response http.ResponseWriter) {
	data, _ := json.Marshal(payload)
	defaults.Logger.Debug("Installation event received")
	defaults.Logger.Debug(string(data))
	response.WriteHeader(http.StatusCreated)
	response.Write(data)
}
