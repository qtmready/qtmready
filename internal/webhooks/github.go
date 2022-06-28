package webhooks

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"

	_zap "go.uber.org/zap"

	"go.breu.io/ctrlplane/internal/defaults"
)

func GithubWebhook(response http.ResponseWriter, request *http.Request) {
	id := request.Header.Get("X-GitHub-Delivery")
	signature := request.Header.Get("X-Hub-Signature")

	if signature == "" {
		handleError(id, ErrorMissingHeaderGithubSignature, http.StatusUnauthorized, response)
		return
	}

	body, _ := ioutil.ReadAll(request.Body)

	if err := verifySignature(body, signature); err != nil {
		handleError(id, err, http.StatusUnauthorized, response)
		return
	}

	headerEvent := request.Header.Get("X-GitHub-Event")

	if headerEvent == "" {
		handleError(id, ErrorMissingHeaderGithubEvent, http.StatusBadRequest, response)
		return
	}

	event := GithubEvent(headerEvent)

	switch event {
	case GithubInstallationEvent:
		var payload GithubInstallationEventPayload
		err := json.Unmarshal(body, &payload)

		if err != nil {
			handleError(id, err, http.StatusBadRequest, response)
			return
		}

		ConsumeGithubInstallationEvent(payload, response)

	default:
		handleError(id, ErrorInvalidEvent, http.StatusBadRequest, response)
	}
}

func verifySignature(payload []byte, signature string) error {
	key := hmac.New(sha1.New, []byte(defaults.Conf.Github.WebhookSecret))
	key.Write(payload)
	result := "sha1=" + hex.EncodeToString(key.Sum(nil))
	defaults.Logger.Debug("ORIG: " + signature)
	defaults.Logger.Debug("RSLT: " + result)
	if result != signature {
		return ErrorVerifySignature
	}
	defaults.Logger.Debug("Signature verified")
	return nil
}

func handleError(requestId string, err error, status int, response http.ResponseWriter) {
	defaults.Logger.Error(err.Error(), _zap.String("request_id", requestId))
	response.WriteHeader(status)
	response.Write([]byte(err.Error()))
}
