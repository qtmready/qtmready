// Crafted with ❤ at Breu, Inc. <info@breu.io>, Copyright © 2024.
//
// Functional Source License, Version 1.1, Apache 2.0 Future License
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
	"time"

	"go.breu.io/quantm/internal/db/fields"
)

type (
	// Int64 is a type alias for int64. Although gocql supports int64, during our application we needed conversions to
	// and from string and int64.
	Int64 = fields.Int64

	// Duration represents a time duration.
	//
	// It addresses shortcomings of Cassandra's duration field, which lacks nanosecond precision.
	Duration = fields.Duration

	// Sensitive represents a string encrypted using AES-GCM.
	//
	// It provides encryption and decryption of sensitive data within the application, ensuring that the data is stored
	// and transmitted securely without being exposed in plain text. This is particularly useful for protecting sensitive
	// values, both at rest and in motion.
	//
	// The encryption key is derived from a secret by calling shared.Service(). shared.Service is a singleton initialized
	// at application startup using environment variables.
	Sensitive = fields.Sensitive
)

func NewDuration(d time.Duration) Duration {
	return fields.NewDuration(d)
}
