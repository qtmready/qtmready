package utils

import (
	"net/http"

	"go.breu.io/ctrlplane/internal/cmn"
)

func HandleHTTPError(writer http.ResponseWriter, err error, status int) { // TODO: get the id from the context
	cmn.Log.Error(err.Error())
	writer.WriteHeader(status)
	writer.Write([]byte(err.Error()))
}

func EncodeJWTPayload(payload map[string]interface{}) string {
	_, out, _ := cmn.JWT.Encode(payload)
	return out
}
