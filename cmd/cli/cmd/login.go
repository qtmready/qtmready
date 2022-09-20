package cmd

import "github.com/spf13/cobra"

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to ctrlplane.ai",
	Long:  "Login to ctrlplane.ai",
}
