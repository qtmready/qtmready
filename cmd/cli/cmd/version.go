// Copyright © 2022, Breu Inc. <info@breu.io>. All rights reserved. 
//
// This software is made available by Breu, Inc., under the terms of the Breu  
// Community License Agreement, Version 1.0 located at  
// http://www.breu.io/breu-community-license/v1. BY INSTALLING, DOWNLOADING,  
// ACCESSING, USING OR DISTRIBUTING ANY OF THE SOFTWARE, YOU AGREE TO THE TERMS  
// OF SUCH LICENSE AGREEMENT. 

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
