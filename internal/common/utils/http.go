package utils

import (
	"net/http"

	"go.breu.io/ctrlplane/internal/common"
)

func HandleHTTPError(writer http.ResponseWriter, err error, status int) { // TODO: get the id from the context
	common.Logger.Error(err.Error())
	writer.WriteHeader(status)
	writer.Write([]byte(err.Error()))
}

func EncodeJWTPayload(payload map[string]interface{}) string {
	_, out, _ := common.JWT.Encode(payload)
	return out
}
