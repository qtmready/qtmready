// Copyright © 2022, Breu Inc. <info@breu.io>. All rights reserved. 

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the current ctrlplane version",
	Long: `
Displays the build id and git short sha for the current binary. This is useful for debugging and reporting issues.
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("TODO: integrate versioning ...") // TODO: integrate versioning
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
