// Copyright © 2023, Breu, Inc. <info@breu.io>. All rights reserved.
//
// This software is made available by Breu, Inc., under the terms of the BREU COMMUNITY LICENSE AGREEMENT, Version 1.0,
// found at https://www.breu.io/license/community. BY INSTALLING, DOWNLOADING, ACCESSING, USING OR DISTRIBUTING ANY OF
// THE SOFTWARE, YOU AGREE TO THE TERMS OF THE LICENSE AGREEMENT.
//
// The above copyright notice and the subsequent license agreement shall be included in all copies or substantial
// portions of the software.
//
// Breu, Inc. HEREBY DISCLAIMS ANY AND ALL WARRANTIES AND CONDITIONS, EXPRESS, IMPLIED, STATUTORY, OR OTHERWISE, AND
// SPECIFICALLY DISCLAIMS ANY WARRANTY OF MERCHANTABILITY OR FITNESS FOR A PARTICULAR PURPOSE, WITH RESPECT TO THE
// SOFTWARE.
//
// Breu, Inc. SHALL NOT BE LIABLE FOR ANY DAMAGES OF ANY KIND, INCLUDING BUT NOT LIMITED TO, LOST PROFITS OR ANY
// CONSEQUENTIAL, SPECIAL, INCIDENTAL, INDIRECT, OR DIRECT DAMAGES, HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY,
// ARISING OUT OF THIS AGREEMENT. THE FOREGOING SHALL APPLY TO THE EXTENT PERMITTED BY APPLICABLE LAW.

package auth

import (
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
