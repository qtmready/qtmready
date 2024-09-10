package mutex

import (
	"errors"
	"fmt"
)

var (
	ErrNilContext   = errors.New("contexts not initialized")
	ErrNoResourceID = errors.New("no resource ID provided")
)

type (
	MutexError struct {
		id   string // the id of the mutex.
		kind string // kind of error. can be "acquire lock", "release lock", or "start workflow".
	}
)

func (e *MutexError) Error() string {
	return fmt.Sprintf("%s: failed to %s.", e.id, e.kind)
}

// NewAcquireLockError creates a new acquire lock error.
func NewAcquireLockError(id string) error {
	return &MutexError{id, "acquire lock"}
}

// NewReleaseLockError creates a new release lock error.
func NewReleaseLockError(id string) error {
	return &MutexError{id, "release lock"}
}

// NewPrepareMutexError creates a new start workflow error.
func NewPrepareMutexError(id string) error {
	return &MutexError{id, "prepare mutex"}
}

func NewCleanupMutexError(id string) error {
	return &MutexError{id, "cleanup mutex"}
}
