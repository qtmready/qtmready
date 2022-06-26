package webhooks

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"go.breu.io/ctrlplane/internal/defaults"
)

func GithubWebhook(response http.ResponseWriter, request *http.Request) {
	// id := request.Header.Get("X-GitHub-Delivery")
	signature := request.Header.Get("X-Hub-Signature")

	if signature == "" {
		defaults.Logger.Error(ErrorMissingHeaderGithubSignature.Error())
		response.WriteHeader(http.StatusUnauthorized)
		response.Write([]byte("Missing X-Hub-Signature Header"))
		return
	}

	body, _ := ioutil.ReadAll(request.Body)

	if err := verifySignature(body, signature); err != nil {
		defaults.Logger.Error(err.Error())
		response.WriteHeader(http.StatusUnauthorized)
		response.Write([]byte("Signature verification failed"))
		return
	}

	headerEvent := request.Header.Get("X-GitHub-Event")

	if headerEvent == "" {
		defaults.Logger.Error(ErrorMissingHeaderGithubEvent.Error())
		response.WriteHeader(http.StatusUnauthorized)
		response.Write([]byte("Missing X-GitHub-Event Header"))
		return
	}

	event := GithubEvent(headerEvent)

	switch event {
	case GithubInstallationEvent:
		var payload GithubInstallationEventPayload
		err := json.Unmarshal(body, &payload)

		if err != nil {
			defaults.Logger.Error(ErrorPayloadParser.Error())
			response.WriteHeader(http.StatusBadRequest)
			response.Write([]byte("Error parsing payload"))
			return
		}
		consumeGithubInstallationEvent(payload, response)
	default:
		defaults.Logger.Error(ErrorInvalidEvent.Error())
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte("Invalid event"))
	}

	// response.Write([]byte("Github Webhook received"))

	// get event from header
	// parse event
	// parse payload
	// verify HMAC
	// handle event
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
