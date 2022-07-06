package models

import (
	"github.com/gosimple/slug"
	"github.com/lithammer/shortuuid/v4"
)

var IsNullUUID = "00000000-0000-0000-0000-000000000000"
var IsNullString = ""

func slugify(s string) string {
	return slug.Make(s) + "-" + shortuuid.New()
}
