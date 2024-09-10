package auth

import (
	"encoding/json"

	"github.com/gocql/gocql"
	"golang.org/x/crypto/bcrypt"

	"go.breu.io/quantm/internal/db"
)

func (t *Team) PreCreate() error { t.Slug = db.CreateSlug(t.Name); return nil }
func (t *Team) PreUpdate() error { return nil }

func (u *User) PreCreate() error { u.SetPassword(u.Password); return nil }
func (u *User) PreUpdate() error { return nil }

func (a *Account) PreCreate() error { return nil }
func (a *Account) PreUpdate() error { return nil }

func (t *TeamUser) PreCreate() error { return nil }
func (t *TeamUser) PreUpdate() error { return nil }

// SetPassword hashes the clear text password using bcrypt.
//
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
//
// TODO: verify that the team exists.
func (u *User) SetActiveTeam(id gocql.UUID) { u.TeamID = id }

// SendVerificationEmail sends a verification email.
func (u *User) SendVerificationEmail() error {
	return nil
}

// SendEmail is the main function responsible for sending emails to users.
func (u *User) SendEmail() error {
	return nil
}

func (mp MessageProviderUserInfo) MarshalCQL(info gocql.TypeInfo) ([]byte, error) {
	return json.Marshal(mp)
}

func (mp *MessageProviderUserInfo) UnmarshalCQL(info gocql.TypeInfo, data []byte) error {
	if len(data) == 0 {
		*mp = MessageProviderUserInfo{}
		return nil
	}

	return json.Unmarshal(data, mp)
}
