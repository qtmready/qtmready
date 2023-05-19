package client

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"go.breu.io/ctrlplane/internal/auth"
	"go.breu.io/ctrlplane/internal/core"
	"go.breu.io/ctrlplane/internal/shared"
)

var Client client

type client struct {
	AuthClient *auth.Client
	CoreClient *core.Client
}

func (c *client) CheckStatus(r *http.Response, successCodes ...int) {
	pass := false
	for _, c := range successCodes {
		if r.StatusCode == c {
			pass = true
		}
	}
	if pass == false {
		fmt.Printf("Command failed with status code: %d\r\n", r.StatusCode)
	}
}

func (c *client) CheckError(err error) {
	if err != nil {
		if strings.Contains(err.Error(), "No connection") {
			fmt.Print("Quantum server is not running\n")
		} else {
			fmt.Printf("Command failed: %v", err.Error())
		}
		os.Exit(1)
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
