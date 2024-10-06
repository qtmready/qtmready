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

package fields

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gocql/gocql"
)

type (
	// Duration represents a time duration.
	//
	// It addresses shortcomings of Cassandra's duration field, which lacks nanosecond precision.
	Duration struct {
		time.Duration
	}
)

// UnmarshalCQL unmarshals a Cassandra duration (stored as text or varchar) into Duration.
//
// Cassandra does not support nanosecond precision in duration fields, so we store durations as strings in the
// database. This function unmarshals a duration string from a Cassandra text or varchar column.
//
// Handles null values.
func (d *Duration) UnmarshalCQL(info gocql.TypeInfo, data []byte) error {
	if info.Type() != gocql.TypeVarchar && info.Type() != gocql.TypeText {
		return fmt.Errorf("expected varchar or text type, got %v", info.Type())
	}

	if len(data) == 0 {
		return nil // Handle null values if needed
	}

	s := string(data)

	duration, err := time.ParseDuration(s)
	if err != nil {
		return fmt.Errorf("error parsing duration %s: %w", s, err)
	}

	d.Duration = duration

	return nil
}

// MarshalCQL marshals Duration into a Cassandra duration.
//
// Converts the Duration to a string and returns it as a byte slice.
func (d Duration) MarshalCQL(info gocql.TypeInfo) ([]byte, error) {
	return []byte(d.String()), nil
}

// MarshalJSON marshals Duration to JSON string.
//
// Converts the Duration to a string and encodes it as a JSON string.
//
// Handles null values.
func (d Duration) MarshalJSON() ([]byte, error) {
	if d.Duration == 0 {
		return []byte("null"), nil // Handle null case for JSON
	}

	return []byte(fmt.Sprintf(`"%s"`, d.String())), nil
}

// UnmarshalJSON unmarshals JSON string into Duration.
//
// Decodes a JSON string into a Duration.
//
// Handles null values.
func (d *Duration) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	if s == "null" {
		return nil // Handle null case for JSON
	}

	duration, err := time.ParseDuration(s)
	if err != nil {
		return fmt.Errorf("error parsing duration %s: %w", s, err)
	}

	d.Duration = duration

	return nil
}

// String returns the string representation of Duration.
//
// Returns the string representation of the Duration.
func (d Duration) String() string {
	return d.Duration.String()
}

func NewDuration(d time.Duration) Duration {
	return Duration{d}
}
