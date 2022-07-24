package models

import (
	"time"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/table"
	"golang.org/x/crypto/bcrypt"

	"go.breu.io/ctrlplane/internal/common"
	"go.breu.io/ctrlplane/internal/db"
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

	teamColumns = []string{
		"id",
		"name",
		"slug",
		"created_at",
		"updated_at",
	}

	teamUsersColumns = []string{
		"id",
		"user_id",
		"team_id",
		"created_at",
		"updated_at",
	}

	userMeta = table.Metadata{
		Name:    "users",
		Columns: userColumns,
	}

	teamMeta = table.Metadata{
		Name:    "teams",
		Columns: teamColumns,
		PartKey: []string{},
		SortKey: []string{},
	}

	teamUsersMeta = table.Metadata{
		Name:    "team_users",
		Columns: teamUsersColumns,
	}

	userTable     = table.New(userMeta)
	teamTable     = table.New(teamMeta)
	teamUserTable = table.New(teamUsersMeta)
)

type (
	Team struct {
		ID        gocql.UUID `json:"id" cql:"id"`
		Name      string     `json:"name" validate:"required"`
		Slug      string     `json:"slug"`
		CreatedAt time.Time  `json:"created_at"`
		UpdatedAt time.Time  `json:"updated_at"`
	}

	User struct {
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

	TeamUser struct {
		ID        gocql.UUID `json:"id" cql:"id"`
		UserID    gocql.UUID `json:"user_id" cql:"user_id"`
		TeamID    gocql.UUID `json:"team_id" cql:"team_id"`
		CreatedAt time.Time  `json:"created_at"`
		UpdatedAt time.Time  `json:"updated_at"`
	}
)

/**
 * Team recievers
 **/

func (t *Team) Save() error {
	if t.ID.String() == NullUUID {
		return t.create()
	} else {
		return t.update()
	}
}

func (t *Team) create() error {
	t.ID, _ = gocql.RandomUUID()
	t.Slug = slugify(t.Name)

	now := time.Now()
	t.CreatedAt = now
	t.UpdatedAt = now

	query := db.DB.Session.Query(teamTable.Insert()).BindStruct(t)

	if err := query.ExecRelease(); err != nil {
		return err
	}

	return nil
}

func (t *Team) update() error {
	t.UpdatedAt = time.Now()

	query := db.DB.Session.Query(teamTable.Update()).BindStruct(t)

	if err := query.ExecRelease(); err != nil {
		return err
	}

	return nil
}

/**
 * User recievers
 **/

func (u *User) Save() error {
	if err := common.Validator.Struct(u); err != nil {
		return err
	}

	if u.ID.String() == NullUUID {
		return u.create()
	} else {
		return u.update()
	}
}

func (u *User) Get(params map[string]interface{}) error {
	query := db.DB.Session.Query(userTable.Get()).BindMap(params)

	if err := query.GetRelease(u); err != nil {
		return err
	}

	return nil
}

func (u *User) Filter(params map[string]interface{}) ([]User, error) {
	users := []User{}
	query := db.DB.Session.Query(userTable.Select()).BindMap(params)

	if err := query.SelectRelease(&users); err != nil {
		return users, err
	}

	return users, nil
}

func (u *User) Team() {}

// Given a password, sets the user's password to a hashed version.
func (u *User) SetPassword(password string) {
	p, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	u.Password = string(p)
}

// Verifies the plain text password against the hashed password.
func (u *User) VerifyPassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) == nil
}

func (u *User) SendEmail() {}
func (u *User) Suspend()   {}
func (u *User) Restore()   {}

// Creates a new user. If email already exists, returns an error.
func (u *User) create() error {
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
func (u *User) update() error {

	u.UpdatedAt = time.Now()

	query := db.DB.Session.Query(userTable.Update()).BindStruct(u)

	if err := query.ExecRelease(); err != nil {
		return err
	}

	return nil
}

/**
 * TeamUser recievers
 **/

func (t *TeamUser) Create() error {
	t.ID, _ = gocql.RandomUUID()

	now := time.Now()
	t.CreatedAt = now
	t.UpdatedAt = now

	query := db.DB.Session.Query(teamUserTable.Insert()).BindStruct(t)

	if err := query.ExecRelease(); err != nil {
		return err
	}

	return nil
}
