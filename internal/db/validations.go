// Copyright © 2022, Breu Inc. <info@breu.io>. All rights reserved.

package db

import (
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2/qb"
	"github.com/scylladb/gocqlx/v2/table"
)

type placeholder struct {
	ID gocql.UUID `json:"id" cql:"id"`
}

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
	var args []reflect.Value // Empty args for reflect.Call

	dest := &placeholder{} // Initializing the empty placeholder to act as a destination for Get call

	tbl := fl.
		Parent().Addr().            // Getting the pointer the parent struct
		MethodByName("GetTable").   // Getting the "GetTable" function by name
		Call(args)[0].              // Calling the function and getting the return value
		Interface().(*table.Table). // Casting the value to *tbl.Table
		Metadata().Name             // Getting the tbl name

	clause := qb.
		EqLit(fl.FieldName(), "'"+fl.Field().Interface().(string)+"'") // forcing args inside '' to provide escaping

	query := qb.
		Select(tbl).                   // Using the qb to compose select query
		Columns("id", fl.FieldName()). // Selecting the return columns
		Where(clause).                 // composing the where clause
		Query(DB.Session)              // using the existing database connection

	err := query.Iter().Unsafe().Get(dest) // Running the "get" query in unsafe mode.

	return err != nil
}
