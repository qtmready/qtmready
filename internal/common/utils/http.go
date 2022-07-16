package utils

import (
	"net/http"

	"go.breu.io/ctrlplane/internal/common"
	"go.uber.org/zap"
)

func HandleHttpError(id string, err error, status int, writer http.ResponseWriter) { // TODO: get the id from the context
	common.Logger.Error(err.Error(), zap.String("request_id", id))
	writer.WriteHeader(status)
	writer.Write([]byte(err.Error()))
}
