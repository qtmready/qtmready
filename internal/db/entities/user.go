package entities

import (
	"time"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/table"
	"golang.org/x/crypto/bcrypt"
)

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

func (u User) GetTable() *table.Table { return userTable }
func (u User) PreCreate() error       { return nil }
func (u User) PreUpdate() error       { return nil }

// Given a password, sets the user's password to a hashed version.
// NOTE: This only updates the field. You will have to run the method to persist the change.
//
//   params := db.QueryParams{"email": "user@example.com"}
//   user, _ := db.Get[User](params)
//   user.SetPassword("password")
//   db.Save[User](user)
func (u *User) SetPassword(password string) {
	p, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	u.Password = string(p)
}

// Verifies the plain text password against the hashed password.
func (u *User) VerifyPassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) == nil
}
