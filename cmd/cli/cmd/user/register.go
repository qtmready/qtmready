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
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	client "go.breu.io/quantm/cmd/cli/apiClient"
	"go.breu.io/quantm/cmd/cli/utils/models"
	"go.breu.io/quantm/internal/auth"
)

type (
	registerOptions struct {
		auth.RegisterationRequest
	}
)

func NewCmdUserRegister() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register",
		Short: "Registers a user",
		Long:  `registers a user with quantum`,

		RunE: func(cmd *cobra.Command, args []string) error {
			f := &registerOptions{}
			if err := tea.NewProgram(models.InitializeInputModel(f)).Start(); err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}

func (o *registerOptions) RunCmd() {
	o.RunRegister()
}

func (o *registerOptions) GetFields() []string {
	fields := []string{"First Name", "Last Name", "Email", "Team", "Password", "Confirm Password"}
	return fields
}

func (o *registerOptions) BindFields(inputs []textinput.Model) {
	for i := 0; i < len(inputs); i++ {
		switch inputs[i].Placeholder {
		case "Email":
			o.RegisterationRequest.Email = inputs[i].Value()
		case "First Name":
			o.RegisterationRequest.FirstName = inputs[i].Value()
		case "Last Name":
			o.RegisterationRequest.LastName = inputs[i].Value()
		case "Team":
			o.RegisterationRequest.TeamName = inputs[i].Value()
		case "Password":
			o.RegisterationRequest.Password = inputs[i].Value()
		case "Confirm Password":
			o.RegisterationRequest.ConfirmPassword = inputs[i].Value()
		}
	}
}

func (o registerOptions) RunRegister() {
	c := client.Client
	r, err := c.AuthClient.Register(context.Background(), o.RegisterationRequest)

	c.CheckError(err)

	pr, err := auth.ParseRegisterResponse(r)
	if err != nil {
		panic("Error: Unable to parse register response")
	}

	switch r.StatusCode {
	case 201:
		fmt.Println("User Registered")
	case 400:
		err, ok := pr.JSON400.Errors.Get("email")
		if ok && err == "already exists" {
			fmt.Println("User already exists")
		} else {
			fmt.Printf("Unable to register user, err:%v\n", err)
		}

	default:
		fmt.Println("Unable to register user")
	}
}

// GetFields gets the fields for prompt from structure tags
// (problem: the names will appear in alphabetical order)
// func (o *registerOptions) GetFields() []string {
// 	t := reflect.TypeOf(o.RegisterationRequest)
// 	fields := make([]string, t.NumField())

// 	for i := 0; i < t.NumField(); i++ {
// 		if val, ok := t.Field(i).Tag.Lookup("json"); ok {
// 			fields[i] = val
// 		}
// 	}

// 	return fields
// }

// BindFields binds the user input with structure
// func (o *registerOptions) BindFields(inputs []textinput.Model) {
// 	v := reflect.ValueOf(o.RegisterationRequest)
// 	t := reflect.TypeOf(o.RegisterationRequest)
// 	for i := 0; i < len(inputs); i++ {
// 		for j := 0; j < t.NumField(); j++ {
// 			tag := t.Field(i).Tag.Get("json")
// 			if tag == inputs[i].Placeholder {
// 				val := v.Field(i)
// 				if val.CanSet() {
// 					val.SetString(inputs[i].Value())
// 				}
// 			}
// 		}
// 	}
// }
