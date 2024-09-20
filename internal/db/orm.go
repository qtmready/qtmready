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


package db

import (
	"reflect"
	"slices"
	"time"

	iqb "github.com/Guilospanck/igocqlx/qb"
	itable "github.com/Guilospanck/igocqlx/table"
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2/qb"
)

type (
	// QueryParams defines the query params required for DB lookup queries.
	QueryParams map[string]string

	// An Entity defines the interface for a database entity.
	Entity interface {
		GetTable() itable.ITable
		PreCreate() error
		PreUpdate() error
	}

	GetOrCreateResponse[T Entity] struct {
		Entity  T
		Created bool
	}
)

func SelectBuilder(table string) iqb.ISelectBuilder {
	return &iqb.SelectBuilder{SB: qb.Select(table)}
}

// Get the entity by given query params.
//
// FIXME: sometimes you have to manually surround the value with "'" to make cql work
//
// A simple example:
//
//	type User struct {
//	  ID     string `json:"getID" cql:"getID"`
//	  Email  string `json:"email" cql:"email"`
//	}
//
//	params := db.QueryParams{"email": "email@example.com"}
//	user := &User{}
//	err := db.Get(user, params)
func Get[T Entity](entity T, params QueryParams) error {
	clause := make([]qb.Cmp, 0)

	for key, value := range params {
		clause = append(clause, qb.EqLit(key, value))
	}

	query := SelectBuilder(entity.GetTable().Name()).
		AllowFiltering().
		Columns(entity.GetTable().Metadata().M.Columns...).
		Where(clause...)

	return DB().Session.Query(query.ToCql()).GetRelease(entity)
}

// Filter the entity by given query params.
//
// FIXME: sometimes you have to manually surround the value with "'" to make cql work
//
// A simple example:
//
//		 type User struct {
//		   ID     string `json:"getID" cql:"getID"`
//		   Email  string `json:"name" cql:"name"`
//		 }
//
//		 params := db.QueryParams{"email": "email@example.com"}
//	   users := make([]User, 0)
//		 err := db.Filter(&User{}, &users, params)
func Filter(entity Entity, dest any, params QueryParams) error {
	clause := make([]qb.Cmp, 0)

	for key, value := range params {
		clause = append(clause, qb.EqLit(key, value))
	}

	query := SelectBuilder(entity.GetTable().Name()).
		AllowFiltering().
		Columns(entity.GetTable().Metadata().M.Columns...).
		Where(clause...)

	return DB().Session.Query(query.ToCql()).SelectRelease(dest)
}

// Save saves the entity. If the entity has an ID, it will be updated. Otherwise,
// it will be created. A pointer to the entity must be passed.
//
//	type User struct {
//	  ID     string `json:"getID" cql:"getID"`
//	  Email  string `json:"name" cql:"name"`
//	}
//
//	user := User{Email: "user@example.com"}
//	user, err := db.Save(&user)
func Save[T Entity](entity T) error {
	pk := _id(entity)

	if pk.String() == NullUUID {
		return Create(entity)
	}

	return Update(entity)
}

// Create creates the entity. The entity value is a pointer to the struct.
func Create[T Entity](entity T) error {
	pk, _ := NewUUID()
	return CreateWithID(entity, pk)
}

// CreateWithID forces the ID while creating. The entity value is a pointer to the struct.
func CreateWithID[T Entity](entity T, pk gocql.UUID) error {
	now := time.Now()

	_set(entity, "ID", pk)
	_set(entity, "CreatedAt", now)
	_set(entity, "UpdatedAt", now)

	if err := entity.PreCreate(); err != nil {
		return err
	}

	query := DB().Session.Query(entity.GetTable().Insert()).BindStruct(entity)

	return query.ExecRelease()
}

// Update updates the entity.
//
// NOTE: The assumption is that ID is the primary key and the first one defined in the struct.
// NOTE: you must pass the complete struct.
func Update[T Entity](entity T) error {
	now := time.Now()
	_set(entity, "UpdatedAt", now)

	tbl := entity.GetTable()
	columns := _delPK(tbl.Metadata().M.Columns, tbl.Metadata().M.PartKey)
	clause := _wherePK(tbl.Metadata().M.PartKey)
	stmnt, names := qb.Update(tbl.Name()).
		Set(columns...).
		Where(clause...).
		ToCql()

	query := DB().Session.Query(stmnt, names)
	query = query.BindStruct(entity)

	return query.ExecRelease()
}

// Delete deletes the row from cassandra. The entity value is a pointer to the struct.
//
// CAUTION: Cassandra has a concept of tombstones. When you delete a row, it is not immediately removed from the database.
// Instead, a tombstone is created, which is a marker that tells Cassandra that the row has been deleted. The tombstone
// will be removed after the gc_grace_seconds period has passed. This is a setting in the table definition.
func Delete[T Entity](entity T) error {
	tbl := entity.GetTable()
	stmnt, names := qb.Delete(tbl.Name()).Where(qb.Eq("id")).ToCql()
	query := DB().Session.Query(stmnt, names)

	return query.BindStruct(entity).ExecRelease()
}

// gets the ID of the entity. The entity value is a pointer to the struct.
func _id(entity Entity) gocql.UUID {
	return reflect.ValueOf(entity).Elem().FieldByName("ID").Interface().(gocql.UUID)
}

// Set the value of the field of the entity. The entity value is a pointer to the struct.
func _set(entity Entity, name string, val any) {
	elem := reflect.ValueOf(entity).Elem()
	elem.FieldByName(name).Set(reflect.ValueOf(val))
}

func _delPK(columns, part []string) []string {
	result := make([]string, 0)

	for _, col := range columns {
		if !slices.Contains(part, col) {
			result = append(result, col)
		}
	}

	return result
}

func _wherePK(keys []string) []qb.Cmp {
	result := make([]qb.Cmp, 0)

	for _, key := range keys {
		result = append(result, qb.Eq(key))
	}

	return result
}
