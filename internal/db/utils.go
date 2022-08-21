package db

import (
	"github.com/gosimple/slug"
	"github.com/lithammer/shortuuid/v4"
)

func CreateSlug(s string) string {
	return slug.Make(s) + "-" + shortuuid.New()
}
