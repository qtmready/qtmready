// Copyright © 2023, Breu, Inc. <info@breu.io>. All rights reserved.
//
// This software is made available by Breu, Inc., under the terms of the BREU COMMUNITY LICENSE AGREEMENT, Version 1.0,
// found at https://www.breu.io/license/community. BY INSTALLING, DOWNLOADING, ACCESSING, USING OR DISTRIBUTING ANY OF
// THE SOFTWARE, YOU AGREE TO THE TERMS OF THE LICENSE AGREEMENT.
//
// The above copyright notice and the subsequent license agreement shall be included in all copies or substantial
// portions of the software.
//
// Breu, Inc. HEREBY DISCLAIMS ANY AND ALL WARRANTIES AND CONDITIONS, EXPRESS, IMPLIED, STATUTORY, OR OTHERWISE, AND
// SPECIFICALLY DISCLAIMS ANY WARRANTY OF MERCHANTABILITY OR FITNESS FOR A PARTICULAR PURPOSE, WITH RESPECT TO THE
// SOFTWARE.
//
// Breu, Inc. SHALL NOT BE LIABLE FOR ANY DAMAGES OF ANY KIND, INCLUDING BUT NOT LIMITED TO, LOST PROFITS OR ANY
// CONSEQUENTIAL, SPECIAL, INCIDENTAL, INDIRECT, OR DIRECT DAMAGES, HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY,
// ARISING OUT OF THIS AGREEMENT. THE FOREGOING SHALL APPLY TO THE EXTENT PERMITTED BY APPLICABLE LAW.

package db

import (
	"reflect"
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

	return DB.Session.Query(query.ToCql()).GetRelease(entity)
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
func Filter(entity Entity, dest interface{}, params QueryParams) error {
	clause := make([]qb.Cmp, 0)

	for key, value := range params {
		clause = append(clause, qb.EqLit(key, value))
	}

	query := SelectBuilder(entity.GetTable().Name()).
		AllowFiltering().
		Columns(entity.GetTable().Metadata().M.Columns...).
		Where(clause...)

	if err := DB.Session.Query(query.ToCql()).SelectRelease(dest); err != nil {
		return err
	}

	return nil
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
	pk := getID(entity)

	if pk.String() == NullUUID {
		return Create(entity)
	}

	return Update(entity)
}

// Create creates the entity. The entity value is a pointer to the struct.
func Create[T Entity](entity T) error {
	pk, _ := gocql.RandomUUID()
	now := time.Now()

	setval(entity, "ID", pk)
	setval(entity, "CreatedAt", now)
	setval(entity, "UpdatedAt", now)

	if err := entity.PreCreate(); err != nil {
		return err
	}

	query := DB.Session.Query(entity.GetTable().Insert()).BindStruct(entity)

	if err := query.ExecRelease(); err != nil {
		return err
	}

	return nil
}

// Update updates the entity.
//
// NOTE: The assumption is that ID is the primary key and the first one defined in the struct.
// NOTE: you must pass the complete struct.
func Update[T Entity](entity T) error {
	now := time.Now()
	setval(entity, "UpdatedAt", now)

	tbl := entity.GetTable()
	columns := tbl.Metadata().M.Columns[1:] // Remove the first element. We are assuming it is the primary key.

	if err := DB.Session.Query(tbl.Update(columns...)).BindStruct(entity).ExecRelease(); err != nil {
		return err
	}

	return nil
}

// gets the ID of the entity. The entity value is a pointer to the struct.
func getID(entity Entity) gocql.UUID {
	return reflect.ValueOf(entity).Elem().FieldByName("ID").Interface().(gocql.UUID)
}

// Set the value of the field of the entity. The entity value is a pointer to the struct.
func setval(entity Entity, name string, val interface{}) {
	elem := reflect.ValueOf(entity).Elem()
	elem.FieldByName(name).Set(reflect.ValueOf(val))
}
