package models

import (
	"net/mail"
	"time"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/table"
	"go.breu.io/ctrlplane/internal/conf"
	"golang.org/x/crypto/bcrypt"
)

var userMeta = table.Metadata{
	Name: "users",
	Columns: []string{
		"id",
		"name",
		"email",
		"password",
		"is_active",
		"is_verified",
		"created_at",
		"updated_at",
	},
}

var userTable = table.New(userMeta)

type User struct {
	ID         gocql.UUID `cql:"id"`
	FirstName  string
	LastName   string
	Email      string
	Password   string
	IsVerified bool
	IsActive   bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// Tries to create a new user. If the email is already in use, returns an error.
func (u *User) Create() error {
	if _, err := mail.ParseAddress(u.Email); err != nil {
		return err
	}

	u.ID, _ = gocql.RandomUUID()

	now := time.Now()
	u.CreatedAt = now
	u.UpdatedAt = now

	query := conf.DB.Session.Query(userTable.Insert()).BindStruct(u)

	if err := query.ExecRelease(); err != nil {
		return err
	}

	return nil
}

func (u *User) Get(params struct{}) error {
	return nil
}

func (u *User) HashPassword(password string) (string, error) {
	p, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(p), nil
}

func (u *User) VerifyPassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) == nil
}

func (u *User) Update()                {}
func (u *User) SendVerificationEmail() {}
func (u *User) Suspend()               {}
func (u *User) Restore()               {}
