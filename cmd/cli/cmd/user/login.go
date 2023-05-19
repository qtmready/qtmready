package user

import (
	"context"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	client "go.breu.io/ctrlplane/cmd/cli/apiClient"
	"go.breu.io/ctrlplane/cmd/cli/utils/models"
	"go.breu.io/ctrlplane/internal/auth"
)

func NewCmdUserLogin() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Logs in a user",
		Long:  `logs in the quantum user`,
		RunE: func(cmd *cobra.Command, args []string) error {
			f := &loginOptions{}
			if err := tea.NewProgram(models.InitializeInputModel(f)).Start(); err != nil {
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

	c := client.Client
	r, err := c.AuthClient.Login(context.Background(), o.LoginRequest)
	c.CheckError(err)
	c.CheckStatus(r, 200)

	pr, err := auth.ParseLoginResponse(r)
	if err != nil {
		panic("Error: User login failed")
	}

	println("User logged in")
	os.Setenv("USER_AUTH_TOKEN", pr.JSON200.AccessToken)
}
