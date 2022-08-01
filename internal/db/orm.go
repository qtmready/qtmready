package db

import (
	"reflect"
	"time"

	"github.com/gocql/gocql"
)

// Get the entity by given query params
func Get[T Entity](params QueryParams) (T, error) {
	entity := *new(T)
	query := DB.Session.Query(entity.GetTable().Get()).BindMap(params)

	if err := query.GetRelease(&entity); err != nil {
		return entity, err
	}
	return entity, nil
}

// Save the entity
func Save[T Entity](entity T) error {
	pk := getpk(entity)

	if pk == NullUUID {
		return create(entity)
	} else {
		return update(entity)
	}
}

func Filter[T Entity](params QueryParams) ([]T, error) {
	return nil, nil
}

func create[T Entity](entity T) error {
	elem := reflect.ValueOf(entity).Elem()
	id, _ := gocql.RandomUUID()
	now := time.Now()

	setvalue(elem, "ID", id)
	setvalue(elem, "CreatedAt", now)
	setvalue(elem, "UpdatedAt", now)

	entity.PreCreate()

	query := DB.Session.Query(entity.GetTable().Insert()).BindStruct(entity)
	if err := query.ExecRelease(); err != nil {
		return err
	}

	return nil
}

func update[T Entity](entity T) error {
	now := time.Now()
	elem := reflect.ValueOf(entity).Elem()
	setvalue(elem, "UpdatedAt", now)

	query := DB.Session.Query(entity.GetTable().Update()).BindStruct(entity)
	if err := query.ExecRelease(); err != nil {
		return err
	}

	return nil
}

func getpk(entity interface{}) string {
	return reflect.ValueOf(entity).Elem().FieldByName("ID").Interface().(string)
}

func setvalue(elem reflect.Value, name string, val any) {
	elem.FieldByName(name).Set(reflect.ValueOf(val))
}
