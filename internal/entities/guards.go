// Copyright © 2022, Breu Inc. <info@breu.io>. All rights reserved.

package entities

import (
	"crypto/rand"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gocql/gocql"
	"github.com/jxskiss/base62"
	"github.com/scylladb/gocqlx/v2/table"
	"golang.org/x/crypto/bcrypt"

	"go.breu.io/ctrlplane/internal/db"
)

var (
	guardColumns = []string{
		"id",
		"name",
		"hashed",
		"lookup_id",
		"lookup_type",
		"created_at",
		"updated_at",
	}

	guardMeta = table.Metadata{
		Name:    "guards",
		Columns: guardColumns,
	}

	guardTable = table.New(guardMeta)
)

// Team is the primary owner of the App & primary driver of system-wide RBAC.
type (
	Guard struct {
		ID         gocql.UUID `json:"id" cql:"id"`
		Name       string     `json:"name" validate:"required"`
		Hashed     string     `json:"hashed" validate:"required"`
		LookupID   gocql.UUID `json:"lookup_id" validate:"required"`
		LookupType string     `json:"lookup_type" validate:"required"`
		CreatedAt  time.Time  `json:"created_at"`
		UpdatedAt  time.Time  `json:"updated_at"`
	}
)

func (g *Guard) GetTable() *table.Table { return guardTable }
func (g *Guard) PreCreate() error       { g.SetHashed(g.Hashed); return nil }
func (g *Guard) PreUpdate() error       { return nil }

// CreatePrefix converts the UUID to a base62 string and use it as the prefix.
func (g *Guard) CreatePrefix() string {
	return base62.EncodeToString(g.LookupID[:])
}

// PrefixToID converts the given prefix to a UUID.
//
// The UUID is used as the lookup_id for the Guard.
func (g *Guard) PrefixToID(prefix string) (gocql.UUID, error) {
	id := gocql.UUID{}
	b, err := base62.DecodeString(prefix)

	if err != nil {
		return id, err
	}

	copy(id[:], b)

	return id, nil
}

// GenerateRandomValue generates a 512 bit random value for the API key.
func (g *Guard) GenerateRandomValue() string {
	bytes := make([]byte, 64) // 64 bytes = 512 bits
	_, _ = rand.Read(bytes)   // Secure random bytes

	return base62.EncodeToString(bytes)
}

func (g *Guard) SetHashed(token string) {
	t, _ := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	g.Hashed = string(t)
}

func (g *Guard) VerifyHashed(token string) bool {
	return bcrypt.CompareHashAndPassword([]byte(g.Hashed), []byte(token)) == nil
}

// ConstructAPIKey constructs a new API Key & return the key plain text and the constructed key.
//
// The plaintext is set in Guard.Hashed and the Guard.PreCreate() function hashes it one save.
func (g *Guard) ConstructAPIKey() (string, string) {
	plaintext := g.GenerateRandomValue()
	key := fmt.Sprintf("%s.%s", g.CreatePrefix(), plaintext)

	return plaintext, key
}

// VerifyAPIKey verifies the API key against the database.
func (g *Guard) VerifyAPIKey(key string) (bool, error) {
	parts := strings.Split(key, ".")
	if len(parts) != 2 {
		return false, errors.New("invalid api key")
	}

	prefix := parts[0]
	hashed := parts[1]
	id, err := g.PrefixToID(prefix)

	if err != nil {
		return false, err
	}

	if err := db.Get(g, db.QueryParams{"lookup_id": id.String()}); err != nil {
		return false, err
	}

	return g.VerifyHashed(hashed), nil
}

// NewForUser creates a new API key for the given user.
//
// NOTE: The Guard.PreCreate() function hashes the plain text value.
func (g *Guard) NewForUser(name string, id gocql.UUID) (string, error) {
	g.Name = name
	g.LookupID = id
	g.LookupType = "user"
	plaintext, key := g.ConstructAPIKey()
	g.Hashed = plaintext

	return key, db.Save(g)
}

// NewForTeam creates a new API key for the given team.
//
// NOTE: One team can have only one API Key.
//
// TODO: Implement unique constraint on lookup_id for team.
func (g *Guard) NewForTeam(id gocql.UUID) (string, error) {
	g.Name = "default"
	g.LookupID = id
	g.LookupType = "team"
	hashed, key := g.ConstructAPIKey()
	g.Hashed = hashed

	return key, db.Save(g)
}
