// Crafted with ❤ at Breu, Inc. <info@breu.io>, Copyright © 2024.
//
// Functional Source License, Version 1.1, Apache 2.0 Future License
//
// We hereby irrevocably grant you an additional license to use the Software under the Apache License, Version 2.0 that
// is effective on the second anniversary of the date we make the Software available. On or after that date, you may use
// the Software under the Apache License, Version 2.0, in which case the following will apply:
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
// the License.
//
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
// an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.

package slack

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/slack-go/slack"

	"go.breu.io/quantm/internal/core/defs"
)

const (
	footer  = "Powered by quantm.io"
	preview = "Message from Quantm" // This is a preview of the message content.
)

func formatLineThresholdExceededAttachment(payload *defs.MessageIOLineExeededPayload) slack.Attachment {
	fields := []slack.AttachmentField{
		createRepositoryField(payload.MessageIOPayload.RepoName, payload.DetectChanges.RepoUrl, true),
		createBranchField(payload.MessageIOPayload.BranchName, payload.DetectChanges.CompareUrl, true),
		{
			Title: "*Threshold*",
			Value: fmt.Sprintf("%d", payload.Threshold),
			Short: true,
		},
		{
			Title: "*Total Lines Count*",
			Value: fmt.Sprintf("%d", payload.DetectChanges.Delta),
			Short: true,
		},
		{
			Title: "*Lines Added*",
			Value: fmt.Sprintf("%d", payload.DetectChanges.Added),
			Short: true,
		},
		{
			Title: "*Lines Deleted*",
			Value: fmt.Sprintf("%d", payload.DetectChanges.Removed),
			Short: true,
		},
		{
			Title: "*Details*",
			Value: fmt.Sprintf("*Number of Files Changed:* %d\n\n*Files Changed*\n%s",
				len(payload.DetectChanges.Modified), formatFilesList(payload.DetectChanges.Modified)),
			Short: false,
		},
	}

	if payload.MessageIOPayload.IsChannel {
		fields = append(fields, createPushedByField(payload.MessageIOPayload.Author, payload.MessageIOPayload.AuthorURL, true))
	}

	return slack.Attachment{
		Color:   "warning",
		Pretext: "The number of lines in this pull request exceeds the allowed threshold. Please review and adjust accordingly.",
		// a plain text summary of the attachment used in clients that don't show formatted text (eg. IRC, mobile notifications).
		Fallback:   preview,
		Fields:     fields,
		MarkdownIn: []string{"fields"},
		Footer:     footer,
		Ts:         json.Number(strconv.FormatInt(time.Now().Unix(), 10)),
	}
}

func formatMergeConflictAttachment(payload *defs.MergeConflictMessage) slack.Attachment {
	fields := []slack.AttachmentField{
		{
			Title: "*Commit SHA*",
			Value: fmt.Sprintf("<%s|%s>", payload.CommitUrl, payload.SHA[:7]),
			Short: true,
		},
		createRepositoryField(payload.MessageIOPayload.RepoName, payload.RepoUrl, true),
		createBranchField(payload.MessageIOPayload.BranchName, payload.CommitUrl, true),
	}

	if payload.MessageIOPayload.IsChannel {
		fields = append(fields, createPushedByField(payload.MessageIOPayload.Author, payload.MessageIOPayload.AuthorURL, true))
	}

	return slack.Attachment{
		Color: "warning",
		Pretext: fmt.Sprintf("A recent commit on defualt branch has caused the merge conflict on <%s|%s> branch.",
			payload.CommitUrl, payload.MessageIOPayload.BranchName),
		// a plain text summary of the attachment used in clients that don't show formatted text (eg. IRC, mobile notifications).
		Fallback:   preview,
		Fields:     fields,
		MarkdownIn: []string{"fields"},
		Footer:     footer,
		Ts:         json.Number(strconv.FormatInt(time.Now().Unix(), 10)),
	}
}

func formatStaleBranchAttachment(payload *defs.MessageIOStaleBranchPayload) slack.Attachment {
	return slack.Attachment{
		Pretext: fmt.Sprintf("Stale branch <%s|%s> is detected on repository <%s|%s>. Please review and take necessary action.",
			payload.CommitUrl, payload.MessageIOPayload.BranchName, payload.RepoUrl, payload.MessageIOPayload.RepoName),
		// a plain text summary of the attachment used in clients that don't show formatted text (eg. IRC, mobile notifications).
		Fallback: preview,
	}
}

func createRepositoryField(repo, url string, short bool) slack.AttachmentField {
	return slack.AttachmentField{
		Title: "*Repository*",
		Value: fmt.Sprintf("<%s|%s>", url, repo),
		Short: short,
	}
}

func createBranchField(branch, url string, short bool) slack.AttachmentField {
	return slack.AttachmentField{
		Title: "*Branch*",
		Value: fmt.Sprintf("<%s|%s>", url, branch),
		Short: short,
	}
}

func createPushedByField(author, url string, short bool) slack.AttachmentField {
	return slack.AttachmentField{
		Title: "*Pushed By*", // TODO - may need to change
		Value: fmt.Sprintf("<%s|%s>", url, author),
		Short: short,
	}
}

func formatFilesList(files []string) string {
	result := ""
	for _, file := range files {
		result += "- " + file + "\n"
	}

	return result
}
