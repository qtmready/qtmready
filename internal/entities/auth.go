package entities

import (
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
