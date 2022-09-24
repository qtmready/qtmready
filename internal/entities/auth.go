// Copyright © 2022, Breu Inc. <info@breu.io>. All rights reserved. 

package entities

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2/table"
	"golang.org/x/crypto/bcrypt"

	"go.breu.io/ctrlplane/internal/db"
)

var (
	teamColumns = []string{
		"id",
		"name",
		"slug",
		"created_at",
		"updated_at",
	}

	teamMeta = table.Metadata{
		Name:    "teams",
		Columns: teamColumns,
	}

	teamTable = table.New(teamMeta)
)

// Team is the primary owner of the App & primary driver of system-wide RBAC.
type Team struct {
	ID        gocql.UUID `json:"id" cql:"id"`
	Name      string     `json:"name" validate:"required"`
	Slug      string     `json:"slug"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (t *Team) GetTable() *table.Table { return teamTable }
func (t *Team) PreCreate() error       { t.Slug = db.CreateSlug(t.Name); return nil }
func (t *Team) PreUpdate() error       { return nil }

var (
	userColumns = []string{
		"id",
		"team_id",
		"first_name",
		"last_name",
		"email",
		"password",
		"is_active",
		"is_verified",
		"created_at",
		"updated_at",
	}

	userMeta = table.Metadata{
		Name:    "users",
		Columns: userColumns,
	}

	userTable = table.New(userMeta)
)

// User defines the auth user. A user can be part of multiple teams. The key User.TeamID represents the default team.
type User struct {
	ID         gocql.UUID `json:"id" cql:"id"`
	TeamID     gocql.UUID `json:"team_id" cql:"team_id"`
	FirstName  string     `json:"first_name"`
	LastName   string     `json:"last_name"`
	Email      string     `json:"email" validate:"email,required,db_unique"`
	Password   string     `json:"-" copier:"-"`
	IsVerified bool       `json:"is_verified"`
	IsActive   bool       `json:"is_active"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

func (u *User) GetTable() *table.Table { return userTable }
func (u *User) PreCreate() error       { u.SetPassword(u.Password); return nil }
func (u *User) PreUpdate() error       { return nil }

// SetPassword hashes the clear text password using bcrypt.
// NOTE: This only updates the field. You will have to run the method to persist the change.
//
//	params := db.QueryParams{"email": "user@example.com"}
//	user, _ := db.Get[User](params)
//	user.SetPassword("password")
//	db.Save(user)
func (u *User) SetPassword(password string) {
	p, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	u.Password = string(p)
}

// VerifyPassword verifies the plain text password against the hashed password.
func (u *User) VerifyPassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) == nil
}

// SetActiveTeam sets the active team for the given user.
// TODO: verify that the team exists
func (u *User) SetActiveTeam(id gocql.UUID) { u.TeamID = id }

// SendVerificationEmail sends a verification email.
func (u *User) SendVerificationEmail() error {
	return nil
}

// SendEmail is the main function responsible for sending emails to users.
func (u *User) SendEmail() error {
	return nil
}

/**
 * TeamUser
 * NOTE: this needs to be implemented
 */

var (
	teamUserColumns = []string{
		"id",
		"user_id",
		"team_id",
		"created_at",
		"updated_at",
	}

	teamUserMeta = table.Metadata{
		Name:    "team_users",
		Columns: teamUserColumns,
	}

	teamUserTable = table.New(teamUserMeta)
)

// TeamUser maintains the relationship between teams and users. One user can be part of multiple teams
// NOTE: this needs to be implemented. The long term plan is that we are going to have relationships and `User.TeamID`
// will represent the primary team. This will be used for the initial setup of the user.
type TeamUser struct {
	ID        gocql.UUID `json:"id" cql:"id"`
	UserID    gocql.UUID `json:"user_id" cql:"user_id"`
	TeamID    gocql.UUID `json:"team_id" cql:"team_id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (tu *TeamUser) GetTable() *table.Table { return teamUserTable }
func (tu *TeamUser) PreCreate() error       { return nil }
func (tu *TeamUser) PreUpdate() error       { return nil }

/**
 * Guards
 **/

var (
	guardColumns = []string{
		"id",
		"name",
		"value",
		"lookup_id",
		"lookup_type",
		"created_at",
		"udpdated_at",
	}

	guardMeta = table.Metadata{
		Name:    "guards",
		Columns: guardColumns,
	}

	guardTable = table.New(guardMeta)
)

// Guard is used to protect the API from unauthorized access.
type Guard struct {
	ID         gocql.UUID `json:"id" cql:"id"`
	Name       string     `json:"name" validate:"required"`
	Hashed     string     `json:"hashed" validate:"required"`
	LookupID   gocql.UUID `json:"lookup_id" validate:"required"`
	LookupType string     `json:"lookup_type" validate:"required"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

func (g *Guard) GetTable() *table.Table { return guardTable }
func (g *Guard) PreCreate() error       { g.SetHashed(g.Hashed); return nil }
func (g *Guard) PreUpdate() error       { return nil }

func (g *Guard) CreatePrefix() string {
	return base64.RawURLEncoding.EncodeToString(g.LookupID[:])
}

func (g *Guard) PrefixToID(prefix string) (gocql.UUID, error) {
	id := gocql.UUID{}
	b, err := base64.RawURLEncoding.DecodeString(prefix)
	if err != nil {
		return id, err
	}

	copy(id[:], b)
	return id, nil
}

func (g *Guard) GenerateRandomValue() string {
	bytes := make([]byte, 64) // 64 bytes = 512 bits
	rand.Read(bytes)          // Secure random bytes
	return base64.RawURLEncoding.EncodeToString(bytes)
}

func (g *Guard) SetHashed(token string) {
	t, _ := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	g.Hashed = string(t)
}

func (g *Guard) VerifyHashed(token string) bool {
	return bcrypt.CompareHashAndPassword([]byte(g.Hashed), []byte(token)) == nil
}

func (g *Guard) ConstructAPIKey() (string, string) {
	value := g.GenerateRandomValue()
	key := fmt.Sprintf("%s.%s", g.CreatePrefix(), value)
	return value, key
}

// VerifyAPIKey verifies the API key against the database.
// TODO: returns the user / team id
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
func (g *Guard) NewForUser(name string, user *User) (string, error) {
	g.Name = name
	g.LookupID = user.ID
	g.LookupType = "user"
	hashed, key := g.ConstructAPIKey()
	g.Hashed = hashed
	return key, db.Save(g)
}

// NewForTeam creates a new API key for the given team.
// NOTE: One team can have only one API Key.
// TODO: Implement unique constraint on lookup_id for team
func (g *Guard) NewForTeam(team *Team) (string, error) {
	g.Name = "default"
	g.LookupID = team.ID
	g.LookupType = "team"
	hashed, key := g.ConstructAPIKey()
	g.Hashed = hashed
	return key, db.Save(g)
}
