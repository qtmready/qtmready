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
