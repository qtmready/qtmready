package models

import (
	"time"

	"github.com/gocql/gocql"
	"github.com/jinzhu/copier"
	"github.com/scylladb/gocqlx/table"
	"golang.org/x/crypto/bcrypt"

	"go.breu.io/ctrlplane/internal/common"
	"go.breu.io/ctrlplane/internal/db"
)

var columns = []string{
	"id",
	"first_name",
	"last_name",
	"email",
	"password",
	"is_active",
	"is_verified",
	"created_at",
	"updated_at",
}

var userMeta = table.Metadata{
	Name:    "users",
	Columns: columns,
}

var userTable = table.New(userMeta)

type User struct {
	ID         gocql.UUID `json:"id" cql:"id"`
	FirstName  string     `json:"first_name"`
	LastName   string     `json:"last_name"`
	Email      string     `json:"email" validate:"email,required,db_unique"`
	Password   string     `json:"-" copier:"-"`
	IsVerified bool       `json:"is_verified"`
	IsActive   bool       `json:"is_active"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// Creates a new user. If email already exists, returns an error.
func (u *User) Create() error {
	if err := common.Validator.Struct(u); err != nil {
		common.Logger.Error(err.Error())
		return err
	}

	u.ID, _ = gocql.RandomUUID()
	u.SetPassword(u.Password)

	now := time.Now()
	u.CreatedAt = now
	u.UpdatedAt = now

	query := db.DB.Session.Query(userTable.Insert()).BindStruct(u)

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

	query := db.DB.Session.Query(userTable.Update()).BindStruct(u)

	if err := query.ExecRelease(); err != nil {
		return err
	}

	return nil
}

// Get a user matching `params`.
// TODO: The input currently is a map[string]interface{}, but it should be a interface{}
func (u *User) Get(params map[string]interface{}) error {
	query := db.DB.Session.Query(userTable.Get()).BindMap(params)

	if err := query.GetRelease(u); err != nil {
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
func (u *User) Save()                  {}
