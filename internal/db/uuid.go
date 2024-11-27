package db

import (
	"github.com/gocql/gocql"
	"github.com/google/uuid"
)

// NewUUID generates a new NewUUID v7 but returns it as a gocql.NewUUID type.
func NewUUID() (gocql.UUID, error) {
	id, _ := uuid.NewV7()
	return gocql.UUIDFromBytes(id[:])
}

// MustUUID generates a new NewUUID v7 but returns it as a gocql.NewUUID type. It panics if an error occurs.
func MustUUID() gocql.UUID {
	id, _ := NewUUID()
	return id
}

// ParseUUID converts a string into a uuid.UUID and returns an error if invalid.
func ParseUUID(input string) (uuid.UUID, error) {
	parsed, err := uuid.Parse(input)
	if err != nil {
		return uuid.Nil, err
	}

	return parsed, nil
}
