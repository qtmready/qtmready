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

package user

import (
	"context"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"go.breu.io/quantm/cmd/cli/api"
	"go.breu.io/quantm/cmd/cli/utils/models"
	"go.breu.io/quantm/internal/auth"
	"go.breu.io/quantm/internal/shared"
)

func NewCmdUserLogin() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Logs in a user",
		Long:  `logs in the quantm user`,
		RunE: func(cmd *cobra.Command, args []string) error {
			f := &loginOptions{}
			if _, err := tea.NewProgram(models.InitializeInputModel(f)).Run(); err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}

type (
	loginOptions struct {
		auth.LoginRequest
	}
)

func (o *loginOptions) RunCmd() {
	o.RunLogin()
}

func (o *loginOptions) GetFields() []string {
	fields := []string{"Email", "Password"}
	return fields
}

func (o *loginOptions) BindFields(inputs []textinput.Model) {
	for i := 0; i < len(inputs); i++ {
		switch inputs[i].Placeholder {
		case "Email":
			o.LoginRequest.Email = inputs[i].Value()
		case "Password":
			o.LoginRequest.Password = inputs[i].Value()
		}
	}
}

func (o loginOptions) RunLogin() {
	c := api.Client
	r, err := c.AuthClient.Login(context.Background(), o.LoginRequest)
	c.CheckError(err)
	c.CheckStatus(r, 200)

	pr, err := auth.ParseLoginResponse(r)
	if err != nil {
		panic("Error: User login failed")
	}

	println("Login Successful")

	err = os.WriteFile(shared.CLI().GetConfigFile(), []byte(pr.JSON200.AccessToken), 0400)
	if err != nil {
		panic(err)
	}
}
