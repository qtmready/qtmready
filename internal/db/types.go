package db

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// ConvertPgTypeUUIDToUUID converts pgtype.UUID to uuid.UUID.
// Returns a zero UUID if the pgtype.UUID is invalid.
func ConvertPgTypeUUIDToUUID(pgUUID pgtype.UUID) uuid.UUID {
	if !pgUUID.Valid {
		return uuid.Nil
	}

	return uuid.Must(uuid.FromBytes(pgUUID.Bytes[:]))
}

// ConvertStringToPgTypeUUID converts a string UUID to pgtype.UUID.
// Returns a zero pgtype.UUID if the string is invalid.
func ConvertStringToPgTypeUUID(id string) pgtype.UUID {
	var pgUUID pgtype.UUID
	if id == "" {
		pgUUID.Valid = false
		return pgUUID
	}

	value, err := uuid.Parse(id)
	if err != nil {
		pgUUID.Valid = false
		return pgUUID
	}

	copy(pgUUID.Bytes[:], value[:])
	pgUUID.Valid = true

	return pgUUID
}

// ConvertBoolToPgTypeBool converts a bool to pgtype.Bool.
func ConvertBoolToPgTypeBool(value bool) pgtype.Bool {
	return pgtype.Bool{
		Bool:  value,
		Valid: true,
	}
}

// Function to convert a string duration to pgtype.Interval.
func IntervalToDurationString(interval pgtype.Interval) string {
	if !interval.Valid {
		return "0s"
	}

	days := time.Duration(interval.Days) * 24 * time.Hour
	micros := time.Duration(interval.Microseconds) * time.Microsecond

	duration := days + micros

	return duration.String()
}

// Function to convert a string duration to pgtype.Interval.
func StringToInterval(durationStr string) pgtype.Interval {
	var interval pgtype.Interval

	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return pgtype.Interval{Valid: false}
	}

	seconds := int64(duration.Seconds())
	microseconds := duration.Microseconds() % 1_000_000
	days := seconds / 86400

	if days < -2147483648 || days > 2147483647 {
		return pgtype.Interval{Valid: false}
	}

	interval.Days = int32(days)
	interval.Microseconds = microseconds
	interval.Valid = true

	return interval
}
