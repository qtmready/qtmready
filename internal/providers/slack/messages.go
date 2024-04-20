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

	"go.breu.io/quantm/internal/core"
)

func FormatLineThresholdExceededMessage(repoName, branchName string, threshold int, branchChanges core.BranchChanges) string {
	actionRequired := "Please review the changes carefully before merging."

	var details string
	if branchChanges.FileCount > 0 {
		details = fmt.Sprintf("This commit introduces significant changes to the codebase.\n\n"+
			"**File Count:** %d\n\n"+
			"**Files:**\n", branchChanges.FileCount)
		for i, file := range branchChanges.Files {
			details += fmt.Sprintf("%d. %s\n", i+1, file)
		}
	} else {
		details = "This commit introduces no changes to the codebase."
	}

	message := fmt.Sprintf(
		":rotating_light: **Early Warning: Line Threshold Exceeded**\n\n"+
			"**Repository:** %s\n"+
			"**Branch:** %s\n"+
			"**Threshold:** %d\n"+
			"**Changes:** %d\n"+
			"**Lines Added:** %d\n"+
			"**Lines Deleted:** %d\n\n"+
			":memo: **Details:**\n"+
			"%s\n\n"+
			":mag: **Action Required:**\n"+
			"%s",
		repoName,
		branchName,
		threshold,
		branchChanges.Changes,
		branchChanges.Additions,
		branchChanges.Deletions,
		details,
		actionRequired,
	)

	return message
}

func FormatMergeConflictMessage(repoName, branchName string) string {
	details :=
		fmt.Sprintf("Merge conflicts were detected when attempting to merge the changes from (%s) into the target branch (main).",
			branchName)
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
			"**Branch:** %s\n"+
			":memo: **Details:**\n\n"+
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
