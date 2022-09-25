// Copyright © 2022, Breu Inc. <info@breu.io>. All rights reserved.

package shared

import (
	"encoding/base64"
	"strings"

	"github.com/gocql/gocql"
)

func CreateGuardPrefix(uuid gocql.UUID) string {
	return base64.RawURLEncoding.EncodeToString([]byte(uuid.String()))
}

func PrefixToID(prefix string) (gocql.UUID, error) {
	b, err := base64.RawURLEncoding.DecodeString(prefix)
	if err != nil {
		return gocql.UUID{}, err
	}

	return gocql.ParseUUID(string(b))
}

func GetPrefixAndToken(key string) (string, string) {
	result := strings.Split(key, ".")
	prefix := result[0]
	token := result[1]

	return prefix, token
}
