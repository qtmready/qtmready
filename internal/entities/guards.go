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
	"go.breu.io/ctrlplane/internal/db"
	"golang.org/x/crypto/bcrypt"
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

// EncodeUUID converts the UUID to a base62 string and use it as the prefix.
func (g *Guard) EncodeUUID(id gocql.UUID) string {
	return base62.EncodeToString(id[:])
}

// DecodeUUID converts the given prefix to a UUID.
func (g *Guard) DecodeUUID(prefix string) (gocql.UUID, error) {
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

// VerifyToken verifies the given api key against the hashed value.
//
// TODO: lookup the relevant parent entity (user or team) first to check if it exists.
func (g *Guard) VerifyToken(token string) bool {
	return bcrypt.CompareHashAndPassword([]byte(g.Hashed), []byte(token)) == nil
}

// ConstructAPIKey constructs a new API Key & return the key plain text and the constructed key.
//
// The plaintext is set in Guard.Hashed and the Guard.PreCreate() function hashes it one save.
func (g *Guard) ConstructAPIKey() (string, string) {
	plaintext := g.GenerateRandomValue()
	key := fmt.Sprintf("%s.%s.%s", g.EncodeUUID(g.ID), g.EncodeUUID(g.LookupID), plaintext)

	return plaintext, key
}

// VerifyAPIKey verifies the API key against the database.
//
// FIXME: fix the lookup, this requires mocking gocqlx.
//
// TODO: implement the database loookup against lookup_id.
//
// TODO: implement the cache so that we don't have to hit the database every time. Possible implementation of good
// key value implementation of LevelDB are:
//
//   - https://github.com/etcd-io/bbolt
//   - https://github.com/dgraph-io/badger
func (g *Guard) VerifyAPIKey(key string) (bool, error) {
	encodedID, encodedLookupID, token, err := g.SplitAPIKey(key)
	if err != nil {
		return false, err
	}

	id, err := g.DecodeUUID(encodedID)
	if err != nil {
		return false, err
	}

	if id != g.ID {
		return false, nil
	}

	lookupID, err := g.DecodeUUID(encodedLookupID)
	if err != nil {
		return false, err
	}

	if lookupID != g.LookupID {
		return false, nil
	}

	return g.VerifyToken(token), nil
}

func (g *Guard) SplitAPIKey(key string) (string, string, string, error) {
	parts := strings.Split(key, ".")
	if len(parts) != 3 {
		return "", "", "", errors.New("invalid api key")
	}

	id := parts[0]
	lookup := parts[1]
	token := parts[2]

	return id, lookup, token, nil
}

// NewForUser creates a new API key for the given user.
//
// NOTE: The Guard.PreCreate() function hashes the plain text value.
func (g *Guard) NewForUser(name string, id gocql.UUID) string {
	g.Name = name
	g.LookupID = id
	g.LookupType = "user"
	plaintext, key := g.ConstructAPIKey()
	g.Hashed = plaintext

	return key
}

// NewForTeam creates a new API key for the given team.
//
// NOTE: One team can have only one API Key.
// TODO: Implement unique constraint on lookup_id for team.
func (g *Guard) NewForTeam(id gocql.UUID) string {
	g.Name = "default"
	g.LookupID = id
	g.LookupType = "team"
	hashed, key := g.ConstructAPIKey()
	g.Hashed = hashed

	return key
}

// GetByEncodedID returns the guard by the given encoded ID.
func (g *Guard) GetByEncodedID(encodedID string) error {
	id, err := g.DecodeUUID(encodedID)
	if err != nil {
		return err
	}

	return db.Get(g, db.QueryParams{"id": id.String()})
}
