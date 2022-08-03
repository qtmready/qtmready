package db

import (
	"reflect"
	"time"

	"github.com/gocql/gocql"
	"go.breu.io/ctrlplane/internal/common"
	"go.uber.org/zap"
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
	common.Logger.Info("Get[T Entity]", zap.Any("params", params))
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

	if pk.String() == NullUUID {
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
	id, err := gocql.RandomUUID()
	if err != nil {
		return err
	}
	now := time.Now()

	setvalue(&entity, "ID", id)
	setvalue(&entity, "CreatedAt", now)
	setvalue(&entity, "UpdatedAt", now)

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
	return reflect.ValueOf(&entity).Elem().Elem().FieldByName("ID").Interface().(gocql.UUID)
}

// Set the value of the field of the entity. The entity value is a pointer to the struct.
func setvalue(entity interface{}, name string, val interface{}) {
	elem := reflect.ValueOf(entity).Elem()
	elem.FieldByName(name).Set(reflect.ValueOf(val))
}
