package fns

import (
	"fmt"
	"strings"
)

// BranchNameFromRef takes a full Git reference string and returns the branch name.
// For example, if the input is "refs/heads/my-branch", the output will be "my-branch".
func BranchNameFromRef(ref string) string {
	return strings.TrimPrefix(ref, "refs/heads/")
}

// BranchNameToRef takes a branch name and returns the full Git reference string.
// For example, if the input is "my-branch", the output will be "refs/heads/my-branch".
func BranchNameToRef(branch string) string {
	return "refs/heads/" + branch
}

func BranchNameToRemoteRef(remote, branch string) string {
	return fmt.Sprintf("refs/remotes/%s/%s", remote, branch)
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
