package db

import (
	"reflect"
	"time"

	"github.com/gocql/gocql"
)

// Get the entity by given query params. A simple example:
//
//   type User struct {
//     ID     string `json:"id" cql:"id"`
//     Email  string `json:"name" cql:"name"`
//   }
//
//   params := db.QueryParams{"email": "email@example.com"}
//   user, err := db.Get[User](params)
func Get[T Entity](params QueryParams) (T, error) {
	entity := *new(T)
	query := DB.Session.Query(entity.GetTable().Get()).BindMap(params)

	if err := query.GetRelease(&entity); err != nil {
		return entity, err
	}
	return entity, nil
}

// Save the entity
//
//   type User struct {
//     ID     string `json:"id" cql:"id"`
//     Email  string `json:"name" cql:"name"`
//   }
//
//   user := User{Email: "user@example.com"}
//   user, err := db.Save[User](user)
func Save[T Entity](entity T) error {
	pk := getid(entity)

	if pk == NullUUID {
		return create(entity)
	} else {
		return update(entity)
	}
}

func Filter[T Entity](params QueryParams) ([]T, error) {
	entity := *new(T)
	var entities []T
	query := DB.Session.Query(entity.GetTable().Select()).BindMap(params)
	if err := query.SelectRelease(&entities); err != nil {
		return entities, err
	}
	return entities, nil
}

func create[T Entity](entity T) error {
	elem := reflect.ValueOf(&entity).Elem() // NOTE: reflect must use pointer
	id, err := gocql.RandomUUID()
	if err != nil {
		return err
	}
	now := time.Now()

	setvalue(elem, "ID", id)
	setvalue(elem, "CreatedAt", now)
	setvalue(elem, "UpdatedAt", now)

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
	elem := reflect.ValueOf(&entity).Elem()
	setvalue(elem, "UpdatedAt", now)

	query := DB.Session.Query(entity.GetTable().Update()).BindStruct(entity)
	if err := query.ExecRelease(); err != nil {
		return err
	}

	return nil
}

func getid(entity Entity) string {
	return reflect.ValueOf(&entity).Elem().FieldByName("ID").Interface().(string)
}

func setvalue(elem reflect.Value, name string, val interface{}) {
	elem.FieldByName(name).Set(reflect.ValueOf(val))
}
