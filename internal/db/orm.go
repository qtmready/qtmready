package db

import (
	"reflect"
	"time"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2/table"
	"go.breu.io/ctrlplane/internal/shared"
)

type (
	// Defines the query params required for DB lookup queries
	QueryParams map[string]interface{}

	// An Entity defines the interface for a database entity
	Entity interface {
		GetTable() *table.Table
		PreCreate() error
		PreUpdate() error
	}
)

// Get the entity by given query params. A simple example:
//
//		type User struct {
//		  ID     string `json:"id" cql:"id"`
//		  Email  string `json:"name" cql:"name"`
//		}
//
//		params := db.QueryParams{"email": "email@example.com"}
//	  user := &User{}
//		err := db.Get(user, params)
func Get[T Entity](entity T, params QueryParams) error {
	query := DB.Session.Query(entity.GetTable().Select()).BindMap(params)

	if err := query.GetRelease(entity); err != nil {
		return err
	}
	return nil
}

// Saves the entity. If the entity has an ID, it will be updated. Otherwise,
// it will be created. A pointer to the entity must be passed.
//
//	type User struct {
//	  ID     string `json:"id" cql:"id"`
//	  Email  string `json:"name" cql:"name"`
//	}
//
//	user := User{Email: "user@example.com"}
//	user, err := db.Save(&user)
func Save[T Entity](entity T) error {
	pk := getID(entity)

	if pk.String() == NullUUID {
		return Create(entity)
	} else {
		return Update(entity)
	}
}

// Filters the entity. NOTE: Work in progress.
func Filter[T any](entity Entity, params QueryParams, columns ...string) ([]T, error) {
	entities := make([]T, 0)
	query := DB.Session.Query(entity.GetTable().Select(columns...)).BindMap(params)
	if err := query.SelectRelease(&entities); err != nil {
		return entities, err
	}
	return entities, nil
}

// Creates the entity. The entity value is a pointer to the struct.
func Create[T Entity](entity T) error {
	id, _ := gocql.RandomUUID()
	now := time.Now()

	setvalue(entity, "ID", id)
	setvalue(entity, "CreatedAt", now)
	setvalue(entity, "UpdatedAt", now)

	if err := entity.PreCreate(); err != nil {
		return err
	}

	query := DB.Session.Query(entity.GetTable().Insert()).BindStruct(entity)
	shared.Logger.Debug("query", "query", query.String())
	if err := query.ExecRelease(); err != nil {
		return err
	}

	return nil
}

// Updates the entity. The assumption is that ID is the primary key and the first one defined in the struct.
func Update[T Entity](entity T) error {
	now := time.Now()
	setvalue(entity, "UpdatedAt", now)

	table := entity.GetTable()
	columns := table.Metadata().Columns[1:] // Remove the first element. We are assuming it is the primary key.

	if err := DB.Session.Query(table.Update(columns...)).BindStruct(entity).ExecRelease(); err != nil {
		return err
	}

	return nil
}

// gets the ID of the entity. The entity value is a pointer to the struct.
func getID(entity Entity) gocql.UUID {
	return reflect.ValueOf(entity).Elem().FieldByName("ID").Interface().(gocql.UUID)
}

// Set the value of the field of the entity. The entity value is a pointer to the struct.
func setvalue(entity Entity, name string, val interface{}) {
	elem := reflect.ValueOf(entity).Elem()
	elem.FieldByName(name).Set(reflect.ValueOf(val))
}
