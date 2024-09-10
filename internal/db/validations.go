// Copyright © 2023, Breu, Inc. <info@breu.io>
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
