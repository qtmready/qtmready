// Copyright © 2022, Breu, Inc. <info@breu.io>. All rights reserved.
//
// This software is made available by Breu, Inc., under the terms of the BREU COMMUNITY LICENSE AGREEMENT, Version 1.0,
// found at https://www.breu.io/license/community. BY INSTALLATING, DOWNLOADING, ACCESSING, USING OR DISTRUBTING ANY OF
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

package db

import (
	"reflect"

	itable "github.com/Guilospanck/igocqlx/table"
	"github.com/go-playground/validator/v10"
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2/qb"
)

type (
	// placeholder is a struct that will be used as a destination for the "get" query during validations.
	placeholder struct {
		ID gocql.UUID `json:"id" cql:"id"`
	}
)

// UniqueField validates that the value of the field is unique in the database.
//
//	type User struct {
//	  Email string `json:"email" validate:"required,email,db_unique"`
//	}
//
//	user := User{Email: "user@example.com"}
//	shared.Validator.Struct(&user)
//
// The validator will check if the field is unique in the database.
func UniqueField(fl validator.FieldLevel) bool {
	var args []reflect.Value // Empty args for reflect.call

	dest := &placeholder{} // Initializing the empty placeholder to act as a destination for Get call

	tbl := fl.
		Parent().Addr().             // Getting the pointer the parent struct
		MethodByName("GetTable").    // Getting the "GetTable" function by name
		Call(args)[0].               // Calling the function and getting the return value
		Interface().(itable.ITable). // Casting the value to *tbl.Table
		Metadata().M.Name            // Getting the tbl name

	clause := qb.
		EqLit(fl.FieldName(), "'"+fl.Field().Interface().(string)+"'") // forcing args inside '' to provide escaping

	query := SelectBuilder(tbl).
		Columns("id", fl.FieldName()).
		Where(clause).
		Query(DB.Session.Session())

	err := query.Iter().Unsafe().Get(dest) // Running the "get" query in unsafe mode.

	return err != nil
}
