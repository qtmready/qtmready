package user

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	client "go.breu.io/ctrlplane/cmd/cli/apiClient"
	"go.breu.io/ctrlplane/cmd/cli/utils/models"
	"go.breu.io/ctrlplane/internal/auth"
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
