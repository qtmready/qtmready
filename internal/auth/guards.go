// Copyright © 2023, Breu, Inc. <info@breu.io>. All rights reserved.
//
// This software is made available by Breu, Inc., under the terms of the BREU COMMUNITY LICENSE AGREEMENT, Version 1.0,
// found at https://www.breu.io/license/community. BY INSTALLING, DOWNLOADING, ACCESSING, USING OR DISTRIBUTING ANY OF
// THE SOFTWARE, YOU AGREE TO THE TERMS OF THE LICENSE AGREEMENT.
//
// The above copyright notice and the subsequent license agreement shall be included in all copies or substantial
// portions of the software.
//
// Breu, Inc. HEREBY DISCLAIMS ANY AND ALL WARRANTIES AND CONDITIONS, EXPRESS, IMPLIED, STATUTORY, OR OTHERWISE, AND
// SPECIFICALLY DISCLAIMS ANY WARRANTY OF MERCHANTABILITY OR FITNESS FOR A PARTICULAR PURPOSE, WITH RESPECT TO THE
// SOFTWARE.
//
// Breu, Inc. SHALL NOT BE LIABLE FOR ANY DAMAGES OF ANY KIND, INCLUDING BUT NOT LIMITED TO, LOST PROFITS OR ANY
// CONSEQUENTIAL, SPECIAL, INCIDENTAL, INDIRECT, OR DIRECT DAMAGES, HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY,
// ARISING OUT OF THIS AGREEMENT. THE FOREGOING SHALL APPLY TO THE EXTENT PERMITTED BY APPLICABLE LAW.

package auth

import (
	"crypto/rand"
	"fmt"
	"strings"
	"time"

	itable "github.com/Guilospanck/igocqlx/table"
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

	guardMeta = itable.Metadata{
		M: &table.Metadata{
			Name:    "guards",
			Columns: guardColumns,
		}}

	guardTable = itable.New(*guardMeta.M)
)

// Team is the primary owner of the Stack & primary driver of system-wide RBAC.
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

func (g *Guard) GetTable() itable.ITable { return guardTable }
func (g *Guard) PreCreate() error        { g.SetHashed(g.Hashed); return nil }
func (g *Guard) PreUpdate() error        { return nil }

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
	bytes := make([]byte, 48)
	_, _ = rand.Read(bytes)

	return base62.EncodeToString(bytes)
}

func (g *Guard) SetHashed(token string) {
	t, _ := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	g.Hashed = string(t)
}

// VerifyToken verifies the given api key against the hashed value.
//
// FIXME: sometimes the bcrypt.CompareHashAndPassword() returns an error even though the token is valid.
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
// TODO: implement the cache so that we don't have to hit the database every time. An in-memory K/V store maybe? We can
// look at some LevelDB's implementations in golang. e.g.
//
//   - https://github.com/etcd-io/bbolt
//   - https://github.com/dgraph-io/badger
func (g *Guard) VerifyAPIKey(key string) error {
	encodedID, encodedLookupID, token, err := g.SplitAPIKey(key)
	if err != nil {
		return err
	}

	id, err := g.DecodeUUID(encodedID)
	if err != nil {
		return ErrMalformedAPIKey
	}

	if err := db.Get(g, db.QueryParams{"id": id.String()}); err != nil {
		return ErrInvalidAPIKey
	}

	if id != g.ID {
		return ErrMalformedAPIKey
	}

	lookupID, err := g.DecodeUUID(encodedLookupID)
	if err != nil {
		return ErrInvalidAPIKey
	}

	if lookupID != g.LookupID {
		return ErrInvalidAPIKey
	}

	if g.VerifyToken(token) {
		return nil
	}

	return ErrCrypto
}

func (g *Guard) SplitAPIKey(key string) (string, string, string, error) {
	parts := strings.Split(key, ".")
	if len(parts) != 3 {
		return "", "", "", ErrMalformedAPIKey
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
	g.ID, _ = gocql.RandomUUID()
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
	g.ID, _ = gocql.RandomUUID()
	g.Name = "default"
	g.LookupID = id
	g.LookupType = "team"
	plaintext, key := g.ConstructAPIKey()
	g.Hashed = plaintext

	return key
}

// Save saves the Guard to the database.
func (g *Guard) Save() error {
	_ = g.PreCreate()
	return db.DB().Session.Query(guardTable.Insert()).BindStruct(g).ExecRelease()
}
