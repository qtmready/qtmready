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

package slack

import (
	"fmt"

	"github.com/slack-go/slack"

	"go.breu.io/quantm/internal/core"
)

func formatLineThresholdExceededAttachment(repoName, branchName string, threshold int, branchChanges core.BranchChanges) slack.Attachment {
	return slack.Attachment{
		Color: "danger",
		Title: "PR Lines Exceed",
		Fields: []slack.AttachmentField{
			createRepositoryField(repoName, branchChanges.RepoUrl),
			createBranchField(branchName, branchChanges.CompareUrl),
			{
				Title: "*Threshold*",
				Value: fmt.Sprintf("%d", threshold),
				Short: true,
			},
			{
				Title: "*Total Lines Count*",
				Value: fmt.Sprintf("%d", branchChanges.Changes),
				Short: true,
			},
			{
				Title: "*Lines Added*",
				Value: fmt.Sprintf("%d", branchChanges.Additions),
				Short: true,
			},
			{
				Title: "*Lines Deleted*",
				Value: fmt.Sprintf("%d", branchChanges.Deletions),
				Short: true,
			},
			{
				Title: "*Details*",
				Value: fmt.Sprintf("*Number of Files Changed:* %d\n*Files Changed:*\n%s",
					branchChanges.FileCount, formatFilesList(branchChanges.Files)),
				Short: false,
			},
		},
		MarkdownIn: []string{"fields"},
	}
}

func formatMergeConflictAttachment(repoName, branchName string) slack.Attachment {
	return slack.Attachment{
		Color: "danger",
		Title: "Merge Conflict",
		Fields: []slack.AttachmentField{
			createRepositoryField(repoName, ""),
			createBranchField(branchName, ""),
		},
		MarkdownIn: []string{"fields"}, // TODO
	}
}

func formatStaleBranchAttachment(repoName, branchName string) slack.Attachment {
	return slack.Attachment{
		Color: "danger",
		Title: "Stale Branch",
		Fields: []slack.AttachmentField{
			createRepositoryField(repoName, ""),
			createBranchField(branchName, ""),
		},
		MarkdownIn: []string{"fields"}, // TODO
	}
}

func createRepositoryField(repoName, repoURL string) slack.AttachmentField {
	return slack.AttachmentField{
		Title: "*Repository*",
		Value: fmt.Sprintf("<%s|%s>", repoURL, repoName),
		Short: true,
	}
}

func createBranchField(branchName, compareUrl string) slack.AttachmentField {
	return slack.AttachmentField{
		Title: "*Branch*",
		Value: fmt.Sprintf("<%s|%s>", compareUrl, branchName),
		Short: true,
	}
}

func formatFilesList(files []string) string {
	result := ""
	for _, file := range files {
		result += "- " + file + "\n"
	}

	return result
}
