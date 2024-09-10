package db

import (
	"reflect"

	itable "github.com/Guilospanck/igocqlx/table"
	"github.com/go-playground/validator/v10"
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2/qb"
)

type (
	// GetPlaceholder is a struct that will be used as a destination for the "get" query during validations.
	GetPlaceholder struct {
		ID gocql.UUID `json:"id" cql:"id"`
	}

	ValidatorFn func(fl validator.FieldLevel) bool
)

func NewGetPlaceholder() *GetPlaceholder {
	return &GetPlaceholder{}
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
	var args []reflect.Value // Empty args for reflect.call

	dest := NewGetPlaceholder() // Initializing the temporary struct to act as a destination for the Get call.

	tbl := fl.
		Parent().Addr().             // Getting the pointer the parent struct
		MethodByName("GetTable").    // Getting the "GetTable" function by name
		Call(args)[0].               // Calling the function and getting the return value
		Interface().(itable.ITable). // Casting the value to *tbl.Table
		Metadata().M.Name            // Getting the tbl name

	clause := qb.
		EqLit(fl.FieldName(), "'"+fl.Field().Interface().(string)+"'") // forcing args inside '' to provide escaping

	query := SelectBuilder(tbl).
		AllowFiltering().
		Columns("id", fl.FieldName()).
		Where(clause)

	err := DB().
		Session.
		Query(query.ToCql()).
		Iter().
		Unsafe().
		Get(dest) // TODO: figure out a way to not use Unsafe()

	return err != nil
}
