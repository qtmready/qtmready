package core

import (
	"fmt"
)

type (
	RepoIORebaseError struct {
		SHA           string
		CommitMessage string
	}
)

func (e *RepoIORebaseError) Error() string {
	return fmt.Sprintf("could not apply %s... %s", e.SHA, e.CommitMessage)
}

func NewRepoIORebaseError(sha, msg string) error {
	return &RepoIORebaseError{
		SHA:           sha,
		CommitMessage: msg,
	}
}
