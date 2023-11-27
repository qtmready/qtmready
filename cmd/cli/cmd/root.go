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

package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"go.breu.io/quantm/cmd/cli/api"
	"go.breu.io/quantm/cmd/cli/cmd/installation"
	"go.breu.io/quantm/cmd/cli/cmd/user"
	"go.breu.io/quantm/internal/shared"
)

var (
	// rootCmd represents the base command when called without any subcommands.
	rootCmd = &cobra.Command{
		Use:   "quantm",
		Short: "quantm is a multi stage release rollout engine with pre-emptive rollbacks.",
		Long: `
quantm is a multi stage release rollout engine for cloud-native applications. It is designed to be used in
conjunction with a CI/CD pipeline & near realtime application monitoring to provide a safe and reliable rollout
process with rollbacks for microservices.

Currently, it only supports monorepo, but poly repo support is on the roadmap.

To learn more, visit https://breu.io/quantm.
  `,
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.AddCommand(NewCmdInit())
	rootCmd.AddCommand(user.NewCmdUser())
	rootCmd.AddCommand(installation.NewCmdInstallation())
	rootCmd.AddCommand(NewCmdVersion())
	rootCmd.AddCommand(NewCmdCDelete())
	rootCmd.AddCommand(NewCmdCGet())
	rootCmd.AddCommand(NewCmdCreate())
	rootCmd.AddCommand(NewCmdEdit())
	rootCmd.AddCommand(NewCmdList())

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	url := shared.CLI().GetURL()
	api.Client.Init(url) // TODO: change to singleton pattern
}
