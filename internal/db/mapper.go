package db

import (
	"github.com/scylladb/go-reflectx"
)

var (
	CQLMapper = reflectx.NewMapperFunc("cql", reflectx.CamelToSnakeASCII)
)
