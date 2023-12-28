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

package git

import (
	"fmt"
	"os"

	goGit "github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
)

func NewCmdMerge() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "merge",
		Short: "Merges the PR of the current branch",
		Long:  `Merges the PR of the current branch`,
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, _ := os.Getwd()
			plainOpenOptions := goGit.PlainOpenOptions{
				DetectDotGit: true,
			}

			repo, err := goGit.PlainOpenWithOptions(dir, &plainOpenOptions)
			if err != nil {
				fmt.Print(err)
				return nil
			}
			currBranch, _ := repo.Head()
			repoConfig, _ := repo.Config()
			remote := repoConfig.Remotes

			fmt.Print(remote["origin"].URLs[0])
			fmt.Printf(currBranch.Name().String())
			return nil
		},
	}

	return cmd
}
