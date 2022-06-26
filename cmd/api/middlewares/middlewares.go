package middlewares

import (
	"net/http"

	_kratos "github.com/ory/kratos-client-go"

	"go.breu.io/ctrlplane/internal/defaults"
)

func KratosMiddleware(next http.Handler) http.Handler {
	conf := _kratos.NewConfiguration()
	conf.Servers = []_kratos.ServerConfiguration{
		{
			URL: defaults.Conf.Kratos.ServerUrl,
		},
	}

	// client := _kratos.NewAPIClient(conf)
	return next
}
