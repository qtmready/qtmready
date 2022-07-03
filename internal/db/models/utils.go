package models

import (
	"github.com/gosimple/slug"
	"github.com/lithammer/shortuuid/v4"
)

func slugify(s string) string {
	return slug.Make(s) + "-" + shortuuid.New()
}
