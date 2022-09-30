// Copyright © 2022, Breu Inc. <info@breu.io>. All rights reserved.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"go.breu.io/ctrlplane/internal/shared"
)

var (
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Show the current ctrlplane version.",
		Long:  `Show the current ctrlplane version.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(shared.Service.Version()) // TODO: integrate versioning
		},
	}
)

func init() {
	rootCmd.AddCommand(versionCmd)
}
