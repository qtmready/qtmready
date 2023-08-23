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
	"errors"
	"fmt"
	"os"
	"regexp"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"

	"go.breu.io/quantm/cmd/cli/utils"

	_ "embed" // required to embed cue files into the binary
)

const (
	dns1035LabelFmt = "[a-z]([-a-z0-9]*[a-z0-9])?"
)

var (
	ErrInvalidLength = errors.New("must be no more than 63 characters")
)

func NewCmdInit() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Create a new quantm project",
		Long: `
Creates a new quantm project in the current directory. This will create .quantm file for
configuration management and .quantm/ directory for state management
  `,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("cmd: Run")
		},
	}

	return cmd
}

func initRun(cmd *cobra.Command, args []string) {
	cwd, err := utils.DetectGitRoot()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	name, err := stackName()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cloud, err := selectCloud()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Selecting sensible defaults for ...", cloud, name)
	fmt.Println("Initializing ...", cwd)
}

func stackName() (string, error) {
	rx := regexp.MustCompile(dns1035LabelFmt) // FIXME: not working!
	prompt := promptui.Prompt{
		Label: "Application Name",
		Validate: func(input string) error {
			if len(input) > 63 {
				return ErrInvalidLength
			}
			if !rx.MatchString(input) {
				return errors.New("must be a valid DNS label")
			}
			return nil
		},
	}

	return prompt.Run()
}

func selectCloud() (string, error) {
	prompt := promptui.Select{
		Label: "Select Cloud Provider",
		Items: []string{"aws", "gcp", "azure"},
	}

	_, result, err := prompt.Run()

	return result, err
}
