package user

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewCmdUserLogin() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Logs in a user",
		Long:  `logs in the quantum user`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("cmd: login")
		},
	}

	return cmd
}
