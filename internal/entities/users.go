// Copyright © 2022, Breu Inc. <info@breu.io>. All rights reserved. 
//
// This software is made available by Breu, Inc., under the terms of the Breu  
// Community License Agreement, Version 1.0 located at  
// http://www.breu.io/breu-community-license/v1. BY INSTALLING, DOWNLOADING,  
// ACCESSING, USING OR DISTRIBUTING ANY OF THE SOFTWARE, YOU AGREE TO THE TERMS  
// OF SUCH LICENSE AGREEMENT. 

package entities

import (
	"time"

	itable "github.com/Guilospanck/igocqlx/table"
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2/table"
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

	userMeta = itable.Metadata{
		M: &table.Metadata{
			Name:    "users",
			Columns: userColumns,
		},
	}

	userTable = itable.New(*userMeta.M)
)

type (
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
)

func (u *User) GetTable() itable.ITable { return userTable }
func (u *User) PreCreate() error        { u.SetPassword(u.Password); return nil }
func (u *User) PreUpdate() error        { return nil }

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
