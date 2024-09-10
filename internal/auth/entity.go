// Crafted with ❤ at Breu, Inc. <info@breu.io>, Copyright © 2022, 2024.
//
// Functional Source License, Version 1.1, Apache 2.0 Future License
//
// We hereby irrevocably grant you an additional license to use the Software under the Apache License, Version 2.0 that
// is effective on the second anniversary of the date we make the Software available. On or after that date, you may use
// the Software under the Apache License, Version 2.0, in which case the following will apply:
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
// the License.
//
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
// an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.


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

// TeamUser returns the team user for the given user.
func (u *User) TeamUser(id gocql.UUID) *TeamUser {
	tu := &TeamUser{}

	if err := db.Get(tu, db.QueryParams{"user_id": u.ID.String(), "team_id": id.String()}); err != nil {
		return nil
	}

	return tu
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
