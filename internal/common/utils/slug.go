package utils

import (
	"github.com/gosimple/slug"
	"github.com/lithammer/shortuuid/v4"
)

var NullUUID = "00000000-0000-0000-0000-000000000000"
var NullString = ""

func CreateSlug(s string) string {
	return slug.Make(s) + "-" + shortuuid.New()
}
