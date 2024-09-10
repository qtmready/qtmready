package db

import (
	"crypto/rand"
	"math/big"
	"strings"

	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
)

var (
	chars = "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func suffix(length int) string {
	sb := strings.Builder{}
	sb.Grow(length)

	for i := 0; i < length; i++ {
		idx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		sb.WriteByte(chars[int(idx.Int64())])
	}

	return sb.String()
}

func CreateSlug(s string) string {
	return slug.Make(s) + "-" + suffix(4)
}

// NewUUID generates a new NewUUID v7 but returns it as a gocql.NewUUID type.
func NewUUID() (gocql.UUID, error) {
	id, _ := uuid.NewV7()
	return gocql.UUIDFromBytes(id[:])
}
