package events

import (
	"github.com/google/uuid"
)

// NewUUID generates a new version 7 UUID.  It returns an error if UUID generation fails.
func NewUUID() (uuid.UUID, error) {
	return uuid.NewV7()
}

// MustUUID generates a new version 7 UUID. It panics if UUID generation fails.
//
// The only condition under which it could theoretically return an error is if the underlying system's source of
// randomness is completely broken or unavailable.  This is an exceptionally rare and serious system-level problem.  It
// would indicate a much deeper issue than just UUID generation. In practice, it almost certainly never going to fail.
// That's why the MustUUID function, which panics on error, is generally considered acceptable in this specific context.
// The panic implies a catastrophic failure of the system's random number generator, which is far more severe than a
// simple UUID generation failure.  A crash due to this problem is arguably preferable to silently generating a
// non-unique or predictable UUID, leading to subtle and hard-to-debug issues.
func MustUUID() uuid.UUID {
	id, err := NewUUID()
	if err != nil {
		panic(err)
	}

	return id
}
