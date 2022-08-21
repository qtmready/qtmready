package entities

import (
	"time"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/table"
	"golang.org/x/crypto/bcrypt"

	"go.breu.io/ctrlplane/internal/db"
)

/**
 * Team
 */

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
		PartKey: []string{},
		SortKey: []string{},
	}

	teamTable = table.New(teamMeta)
)

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

func (t *Team) Users() ([]User, error) {
	users, err := db.Filter[User](&User{}, db.QueryParams{"team_id": t.ID})
	if err != nil {
		return []User{}, err
	}
	return users, nil
}

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

/**
 * User
 */

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

// Given a password, sets the user's password to a hashed version.
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

// Verifies the plain text password against the hashed password.
func (u *User) VerifyPassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) == nil
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
