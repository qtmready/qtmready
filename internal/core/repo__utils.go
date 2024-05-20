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

package core

import (
	"strings"
)

// BranchNameFromRef takes a full Git reference string and returns the branch name.
// For example, if the input is "refs/heads/my-branch", the output will be "my-branch".
func BranchNameFromRef(ref string) string {
	return strings.TrimPrefix(ref, "refs/heads/")
}

// RefFromBranchName takes a branch name and returns the full Git reference string.
// For example, if the input is "my-branch", the output will be "refs/heads/my-branch".
func RefFromBranchName(branch string) string {
	return "refs/heads/" + branch
}

// CreateQuantmRef takes a branch name and returns the full Git reference string for a quantum branch.
// For example, if the input is "my-branch", the output will be "refs/heads/quantm/my-branch".
func CreateQuantmRef(branch string) string {
	return "refs/heads/qtm/" + branch
}

// IsQuantmRef checks if a given Git reference string is a quantum branch reference.
// It returns true if the reference starts with "refs/heads/quantm/", otherwise false.
func IsQuantmRef(ref string) bool {
	return strings.HasPrefix(ref, "refs/heads/qtm/")
}

// IsQuantmBranch returns true if the given branch name starts with "qtm/".
// This is a helper function used to identify branches that are part of the Quantm project.
func IsQuantmBranch(branch string) bool {
	return strings.HasPrefix(branch, "qtm/")
}
