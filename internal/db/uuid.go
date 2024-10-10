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
