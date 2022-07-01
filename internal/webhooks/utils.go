package webhooks

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"net/http"

	"go.uber.org/zap"

	"go.breu.io/ctrlplane/internal/conf"
)

// VerifySignature verifies the signature of a request.
func verifySignature(payload []byte, signature string) error {
	key := hmac.New(sha1.New, []byte(conf.Github.WebhookSecret))
	key.Write(payload)
	result := "sha1=" + hex.EncodeToString(key.Sum(nil))
	conf.Logger.Debug("ORIG: " + signature)
	conf.Logger.Debug("RSLT: " + result)
	if result != signature {
		return ErrorVerifySignature
	}
	conf.Logger.Debug("Signature verified")
	return nil
}

// handleError handles an error and writes it to the response.
func handleError(id string, err error, status int, response http.ResponseWriter) {
	conf.Logger.Error(err.Error(), zap.String("request_id", id))
	response.WriteHeader(status)
	response.Write([]byte(err.Error()))
}
