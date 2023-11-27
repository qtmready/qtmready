// Copyright © 2023, Breu, Inc. <info@breu.io>. All rights reserved.
//
// This software is made available by Breu, Inc., under the terms of the BREU COMMUNITY LICENSE AGREEMENT, Version 1.0,
// found at https://www.breu.io/license/community. BY INSTALLING, DOWNLOADING, ACCESSING, USING OR DISTRIBUTING ANY OF
// THE SOFTWARE, YOU AGREE TO THE TERMS OF THE LICENSE AGREEMENT.
//
// The above copyright notice and the subsequent license agreement shall be included in all copies or substantial
// portions of the software.
//
// Breu, Inc. HEREBY DISCLAIMS ANY AND ALL WARRANTIES AND CONDITIONS, EXPRESS, IMPLIED, STATUTORY, OR OTHERWISE, AND
// SPECIFICALLY DISCLAIMS ANY WARRANTY OF MERCHANTABILITY OR FITNESS FOR A PARTICULAR PURPOSE, WITH RESPECT TO THE
// SOFTWARE.
//
// Breu, Inc. SHALL NOT BE LIABLE FOR ANY DAMAGES OF ANY KIND, INCLUDING BUT NOT LIMITED TO, LOST PROFITS OR ANY
// CONSEQUENTIAL, SPECIAL, INCIDENTAL, INDIRECT, OR DIRECT DAMAGES, HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY,
// ARISING OUT OF THIS AGREEMENT. THE FOREGOING SHALL APPLY TO THE EXTENT PERMITTED BY APPLICABLE LAW.

package installation

import (
	"context"
	"net/http"
	"os"

	"github.com/spf13/cobra"

	"go.breu.io/quantm/cmd/cli/api"
	"go.breu.io/quantm/internal/providers/github"
	"go.breu.io/quantm/internal/shared"
)

func NewCmdInstallationComplete() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "complete",
		Short: "Completes github app installation",
		Long:  `Completes github app installation`,
		Run: func(cmd *cobra.Command, args []string) {
			CompleteInstallation(cmd)
		},
	}

	cmd.Flags().Int64("installation_id", 0, "give installation id of the github app")

	return cmd
}

func CompleteInstallation(cmd *cobra.Command) {
	id, _ := cmd.Flags().GetInt64("installation_id")

	completeInstallationBody := github.CompleteInstallationRequest{
		InstallationID: id,
		SetupAction:    github.SetupActionCreated,
	}

	c := api.Client
	r, err := c.GithubClient.GithubCompleteInstallation(context.Background(), completeInstallationBody, AddAuthHeader)

	defer func() { _ = r.Body.Close() }()
	c.CheckError(err)
	c.CheckStatus(r, 200)

	println("Github app installation complete")
}

// AddAuthHeader adds the authorization header to the request
//
// TODO: get the file path as an environment variable.
func AddAuthHeader(ctx context.Context, req *http.Request) error {
	b, err := os.ReadFile(shared.CLI().GetConfigFile())
	if err != nil {
		panic(err)
	}

	token := string(b)
	req.Header.Set("authorization", "Token "+token)

	return nil
}
