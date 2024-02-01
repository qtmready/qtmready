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
	"context"
	"fmt"
	"os"
	"strings"

	goGit "github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"

	"go.breu.io/quantm/cmd/cli/api"
	"go.breu.io/quantm/internal/providers/github"
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
			headRef, _ := repo.Head()
			currBranch := headRef.Name().Short()
			repoConfig, _ := repo.Config()
			remote := repoConfig.Remotes

			splitURL := func(url string) []string {
				parts := strings.Split(strings.TrimSuffix(strings.TrimPrefix(url, "https://github.com/"), ".git"), "/")
				return parts
			}

			url := remote["origin"].URLs[0]
			urlParts := splitURL(url)
			repoOwner := urlParts[0]
			repoName := urlParts[1]

			fmt.Println(repoName)
			fmt.Println(repoOwner)
			fmt.Println(currBranch)

			mergeDetails := github.CliGitMerge{
				Branch:    currBranch,
				RepoName:  repoName,
				RepoOwner: repoOwner,
			}
			c := api.Client
			req, err := c.GithubClient.CliGitMerge(context.Background(), mergeDetails)
			if err != nil {
				fmt.Print(err.Error())
			}

			fmt.Println(req)
			return nil
		},
	}

	return cmd
}
