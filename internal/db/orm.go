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
	"log/slog"
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

	// Entity defines the interface for a database entity.
	Entity interface {
		GetTable() itable.ITable
		PreCreate() error
		PreUpdate() error
	}

	// GetOrCreateResponse represents the response for GetOrCreate operations.
	GetOrCreateResponse[T Entity] struct {
		Entity  T
		Created bool
	}
)

// SelectBuilder returns a new SelectBuilder for the given table name.
func SelectBuilder(table string) iqb.ISelectBuilder {
	return &iqb.SelectBuilder{SB: qb.Select(table)}
}

// Get retrieves the entity matching the given query parameters.
//
// The function iterates over the provided query parameters and constructs an `EqLit` clause for each key-value pair.
// It then builds a `Select` query using the `SelectBuilder` and executes it against the database.
//
// The query uses `AllowFiltering()` to enable filtering on non-primary key columns.
//
// The results are then mapped to the provided entity using `GetRelease()`.
//
// NOTE: In some cases, manual surrounding of values with single quotes (”) might be required for CQL to work correctly.
//
// Example:
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

	stmnt, names := query.ToCql()
	slog.Info("Get", "stmt", stmnt, "names", names)

	return DB().Session.Query(query.ToCql()).GetRelease(entity)
}

// Filter retrieves all entities matching the given query parameters.
//
// The function iterates over the provided query parameters and constructs an `EqLit` clause for each key-value pair.
// It then builds a `Select` query using the `SelectBuilder` and executes it against the database.
//
// The query uses `AllowFiltering()` to enable filtering on non-primary key columns.
//
// The results are then mapped to the provided destination using `SelectRelease()`.
//
// NOTE: In some cases, manual surrounding of values with single quotes (”) might be required for CQL to work correctly.
//
// Example:
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

// Save saves the entity to the database.
//
// If the entity has an ID, it will be updated. Otherwise, a new entity will be created.
//
// Example:
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

// Create creates a new entity in the database.
//
// The function generates a new UUID for the entity's ID, sets the `CreatedAt` and `UpdatedAt` fields to the current time,
// and calls `PreCreate()` on the entity to allow for pre-creation logic.
//
// It then constructs an `Insert` query and executes it against the database.
//
// Example:
//
//	type User struct {
//	  ID     string `json:"getID" cql:"getID"`
//	  Email  string `json:"name" cql:"name"`
//	}
//
//	user := User{Email: "user@example.com"}
//	err := db.Create(&user)
func Create[T Entity](entity T) error {
	now := time.Now()

	pk, _ := NewUUID()

	_set(entity, "CreatedAt", now)
	_set(entity, "UpdatedAt", now)

	return CreateWithID(entity, pk)
}

// CreateWithID creates a new entity in the database with a specified ID.
//
// The function sets the `CreatedAt` and `UpdatedAt` fields to the current time and calls `PreCreate()` on the entity
// to allow for pre-creation logic. It then constructs an `Insert` query and executes it against the database.
//
// Example:
//
//	type User struct {
//	  ID     string `json:"getID" cql:"getID"`
//	  Email  string `json:"name" cql:"name"`
//	}
//
//	user := User{Email: "user@example.com"}
//	user.ID, _ = db.NewUUID()
//	err := db.CreateWithID(&user, user.ID)
func CreateWithID[T Entity](entity T, pk gocql.UUID) error {
	_set(entity, "ID", pk)

	if err := entity.PreCreate(); err != nil {
		return err
	}

	query := DB().Session.Query(entity.GetTable().Insert()).BindStruct(entity)

	return query.ExecRelease()
}

// Update updates an existing entity in the database.
//
// The function sets the `UpdatedAt` field to the current time and constructs an `Update` query with the entity's
// primary key as the `Where` clause. It then executes the query against the database.
//
// NOTE: The assumption is that the `ID` field is the primary key and the first one defined in the struct.
//
// Example:
//
//	type User struct {
//	  ID     string `json:"getID" cql:"getID"`
//	  Email  string `json:"name" cql:"name"`
//	}
//
//	user := User{ID: "some_id", Email: "updated_user@example.com"}
//	err := db.Update(&user)
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

// Delete removes an entity from the database.
//
// The function constructs a `Delete` query with the entity's primary key as the `Where` clause and executes it
// against the database.
//
// NOTE: Cassandra uses tombstones for deletion, which are markers that indicate a row has been deleted. The tombstone
// will be removed after the `gc_grace_seconds` period has passed, as defined in the table definition.
//
// Example:
//
//	type User struct {
//	  ID     string `json:"getID" cql:"getID"`
//	  Email  string `json:"name" cql:"name"`
//	}
//
//	user := User{ID: "some_id"}
//	err := db.Delete(&user)
func Delete[T Entity](entity T) error {
	tbl := entity.GetTable()
	stmnt, names := qb.Delete(tbl.Name()).Where(qb.Eq("id")).ToCql()
	query := DB().Session.Query(stmnt, names)

	return query.BindStruct(entity).ExecRelease()
}

// _id retrieves the `ID` field of the entity.
func _id(entity Entity) gocql.UUID {
	return reflect.ValueOf(entity).Elem().FieldByName("ID").Interface().(gocql.UUID)
}

// _set sets the value of the field of the entity.
func _set(entity Entity, name string, val any) {
	elem := reflect.ValueOf(entity).Elem()
	elem.FieldByName(name).Set(reflect.ValueOf(val))
}

// _delPK removes the primary key columns from the list of columns.
func _delPK(columns, part []string) []string {
	result := make([]string, 0)

	for _, col := range columns {
		if !slices.Contains(part, col) {
			result = append(result, col)
		}
	}

	return result
}

// _wherePK constructs a `Where` clause for the primary key columns.
func _wherePK(keys []string) []qb.Cmp {
	result := make([]qb.Cmp, 0)

	for _, key := range keys {
		result = append(result, qb.Eq(key))
	}

	return result
}
