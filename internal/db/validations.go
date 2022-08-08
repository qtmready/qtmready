package db

import (
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/table"
	"github.com/scylladb/gocqlx/v2/qb"
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
//	cmn.Validator.Struct(&user)
//
// The validator will check if the field is unique in the database.
func UniqueField(fl validator.FieldLevel) bool {
	args := []reflect.Value{} // Empty args for relect.Call
	dest := &placeholder{}    // Initializing the empty placeholder to act as a destination for Get call

	table := fl.
		Parent().Addr().            // Getting the pointer the parent struct
		MethodByName("GetTable").   // Getting the "GetTable" function by name
		Call(args)[0].              // Calling the function and getting the return value
		Interface().(*table.Table). // Casting the value to *table.Table
		Metadata().Name             // Getting the table name

	clause := qb.
		EqLit(fl.FieldName(), "'"+fl.Field().Interface().(string)+"'") // forcing args inside '' to provide escaping

	query := qb.
		Select(table).                 // Using the querybuilder to compose select query
		Columns("id", fl.FieldName()). // Selecting the return columns
		Where(clause).                 // composing the where clause
		Query(DB.Session)              // using the existing database connection

	err := query.Iter().Unsafe().Get(dest) // Running the "get" query in unsafe mode.
	return err != nil
}
