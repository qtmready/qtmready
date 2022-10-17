// Copyright © 2022, Breu Inc. <info@breu.io>. All rights reserved.

package cmd

import (
	_ "embed" // required to embed cue files into the binary
	"errors"
	"fmt"
	"os"
	"regexp"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"

	"go.breu.io/ctrlplane/cmd/cli/utils"
)

const (
	dns1035LabelFmt = "[a-z]([-a-z0-9]*[a-z0-9])?"
)

var (
	ErrInvalidLength = errors.New("must be no more than 63 characters")

	initCmd = &cobra.Command{
		Use:   "init",
		Short: "Create a new ctrlplane project",
		Long: `
Creates a new ctrlplane project in the current directory. This will create .ctrlplane file for
configuration management and .ctrlplane/ directory for state management
  `,
		Run: initRun,
	}
)

func init() {
	rootCmd.AddCommand(initCmd)
}

func initRun(cmd *cobra.Command, args []string) {
	cwd, err := utils.DetectGitRoot()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// fmt.Println(r)
	name, err := promptName()
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

func promptName() (string, error) {
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
