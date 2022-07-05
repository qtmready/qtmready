package utils

import (
	"net/http"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v45/github"
	"go.breu.io/ctrlplane/internal/conf"
)

func GithubClient(installationID int64) (*github.Client, error) {
	transport, err := ghinstallation.New(http.DefaultTransport, conf.Github.AppID, installationID, []byte(conf.Github.PrivateKey))

	if err != nil {
		return nil, err
	}

	client := github.NewClient(&http.Client{Transport: transport})
	return client, nil
}
