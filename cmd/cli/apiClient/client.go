package client

import (
	"net/http"
	"os"

	"go.breu.io/ctrlplane/internal/auth"
	"go.breu.io/ctrlplane/internal/core"
	"go.breu.io/ctrlplane/internal/shared"
)

var Client client

type client struct {
	AuthClient *auth.Client
	CoreClient *core.Client
}

func checkHttpError(r *http.Response, successCodes ...int) {
	for c := range successCodes {
		if r.StatusCode != c {

			os.Exit(3)
		}
	}
}

// init initializes the auth and core clients to connect with quantum
func (c *client) Init() {

	var err error
	c.AuthClient, err = auth.NewClient(shared.Service.CLI.BaseURL)
	if err != nil {
		c.AuthClient = nil
	}

	c.CoreClient, err = core.NewClient(shared.Service.CLI.APIKEY)
	if err != nil {
		c.CoreClient = nil
	}
}
