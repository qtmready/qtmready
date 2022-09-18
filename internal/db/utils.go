package db

import (
	"strings"

	"github.com/gosimple/slug"
	"github.com/lithammer/shortuuid/v4"
)

func CreateSlug(s string) string {
	return slug.Make(s) + "-" + strings.ToLower(shortuuid.New())
}
