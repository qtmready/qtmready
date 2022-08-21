package db

import (
	"reflect"
	"time"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/table"
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
	query := DB.Session.Query(entity.GetTable().Get()).BindMap(params)

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
	pk := getid(entity)

	if pk.String() == NullUUID {
		return create(entity)
	} else {
		return update(entity)
	}
}

// Filters the entity. NOTE: Work in progress.
func Filter[OUT any](entity Entity, params QueryParams, columns ...string) ([]OUT, error) {
	entities := make([]OUT, 0)
	query := DB.Session.Query(entity.GetTable().Select(columns...)).BindMap(params)
	if err := query.SelectRelease(&entities); err != nil {
		return entities, err
	}
	return entities, nil
}

func create[T Entity](entity T) error {
	id, err := gocql.RandomUUID()
	if err != nil {
		return err
	}
	now := time.Now()

	setvalue(entity, "ID", id)
	setvalue(entity, "CreatedAt", now)
	setvalue(entity, "UpdatedAt", now)

	if err := entity.PreCreate(); err != nil {
		return err
	}

	query := DB.Session.Query(entity.GetTable().Insert()).BindStruct(entity)
	if err := query.ExecRelease(); err != nil {
		return err
	}

	return nil
}

func update[T Entity](entity T) error {
	now := time.Now()
	setvalue(entity, "UpdatedAt", now)

	query := DB.Session.Query(entity.GetTable().Update()).BindStruct(entity)
	if err := query.ExecRelease(); err != nil {
		return err
	}

	return nil
}

func getid(entity Entity) gocql.UUID {
	return reflect.ValueOf(entity).Elem().FieldByName("ID").Interface().(gocql.UUID)
}

// Set the value of the field of the entity. The entity value is a pointer to the struct.
func setvalue(entity interface{}, name string, val interface{}) {
	elem := reflect.ValueOf(entity).Elem()
	elem.FieldByName(name).Set(reflect.ValueOf(val))
}
