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
)

func FormatLineThresholdExceededMessage(repoName, branchName string, changes, threshold int) string {
	details := "This commit introduces significant changes to the codebase."
	actionRequired := "Please review the changes carefully before merging."

	message := fmt.Sprintf(
		":rotating_light: **Early Warning: Line Threshold Exceeded**\n\n"+
			"**Repository:** %s\n"+
			"**Branch:** %s\n"+
			"**Lines Added:** %d\n"+
			"**Threshold:** %d\n\n"+
			":memo: **Details:**\n"+
			"%s\n\n"+
			":mag: **Action Required:**\n"+
			"%s",
		repoName,
		branchName,
		changes,
		threshold,
		details,
		actionRequired,
	)

	return message
}
func FormatMergeConflictMessage(repoName, branchName string) string {
	details := "Merge conflicts were detected when attempting to merge the changes from your branch into the target branch."
	actionRequired := "Please resolve these conflicts before attempting to merge again."

	message := fmt.Sprintf(
		":rotating_light: **Early Warning: Merge Conflicts are expected on branch**\n\n"+
			"**Repository:** %s\n"+
			"**Branch:** %s\n\n"+
			":memo: **Details:**\n"+
			"%s\n\n"+
			":mag: **Action Required:**\n"+
			"%s",
		repoName,
		branchName,
		details,
		actionRequired,
	)

	return message
}

func FormatStaleBranchMessage(repoName, branchName string) string {
	details := "This branch has not had any recent activity."
	actionRequired := "Please check if this branch is still needed. If not, consider deleting it to keep the repository clean."

	message := fmt.Sprintf(
		":rotating_light: **Early Warning: Stale Branch Detected**\n\n"+
			"**Repository:** %s\n"+
			"**Branch:** %s\n\n"+
			":memo: **Details:**\n"+
			"%s\n\n"+
			":mag: **Action Required:**\n"+
			"%s",
		repoName,
		branchName,
		details,
		actionRequired,
	)

	return message
}
