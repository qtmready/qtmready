// Copyright © 2023, Breu, Inc. <info@breu.io>
//
// We hereby irrevocably grant you an additional license to use the Software under the Apache License, Version 2.0 that
// is effective on the second anniversary of the date we make the Software available. On or after that date, you may use
// the Software under the Apache License, Version 2.0, in which case the following will apply:
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
// the License.
//
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
// an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.

package db

import (
	"crypto/rand"
	"math/big"
	"strings"

	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
)

var (
	chars = "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func suffix(length int) string {
	sb := strings.Builder{}
	sb.Grow(length)

	for i := 0; i < length; i++ {
		idx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		sb.WriteByte(chars[int(idx.Int64())])
	}

	return sb.String()
}

func CreateSlug(s string) string {
	return slug.Make(s) + "-" + suffix(4)
}

// NewUUID generates a new NewUUID v7 but returns it as a gocql.NewUUID type.
func NewUUID() (gocql.UUID, error) {
	id, _ := uuid.NewV7()
	return gocql.UUIDFromBytes(id[:])
}
