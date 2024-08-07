package code

import (
	"fmt"
)

type (
	RebaseError struct {
		SHA           string
		CommitMessage string
	}
)

func (e *RebaseError) Error() string {
	return fmt.Sprintf("could not apply %s... %s", e.SHA, e.CommitMessage)
}

func NewRebaseError(sha, msg string) error {
	return &RebaseError{
		SHA:           sha,
		CommitMessage: msg,
	}
}
