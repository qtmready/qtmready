package models

import (
	"net/mail"
	"time"

	"github.com/gocql/gocql"
	"github.com/jinzhu/copier"
	"github.com/scylladb/gocqlx/table"
	"golang.org/x/crypto/bcrypt"

	"go.breu.io/ctrlplane/internal/conf"
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
	Password   string `copier:"-"`
	IsVerified bool
	IsActive   bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// Creates a new user. If email already exists, returns an error.
func (u *User) Create(params interface{}) error {
	if _, err := mail.ParseAddress(u.Email); err != nil {
		return err
	}

	u.ID, _ = gocql.RandomUUID()
	u.SetPassword(u.Password)

	now := time.Now()
	u.CreatedAt = now
	u.UpdatedAt = now

	query := conf.DB.Session.Query(userTable.Insert()).BindStruct(u)

	if err := query.ExecRelease(); err != nil {
		return err
	}

	return nil
}

// Updates a user.
func (u *User) Update(params interface{}) error {
	if err := copier.Copy(u, params); err != nil {
		return err
	}

	u.UpdatedAt = time.Now()

	query := conf.DB.Session.Query(userTable.Update()).BindStruct(u)

	if err := query.ExecRelease(); err != nil {
		return err
	}

	return nil
}

// Get a user matching `params`.
func (u *User) Get(params struct{}) error {
	query := conf.DB.Session.Query(githubInstallationTable.Select()).BindStruct(params)

	if err := query.GetRelease(&u); err != nil {
		return err
	}

	return nil
}

// Given a password, sets the user's password to a hashed version.
func (u *User) SetPassword(password string) {
	p, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	u.Password = string(p)
}

// Verifies the plain text password against the hashed password.
func (u *User) VerifyPassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) == nil
}

func (u *User) SendVerificationEmail() {}
func (u *User) Suspend()               {}
func (u *User) Restore()               {}
