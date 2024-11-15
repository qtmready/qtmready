package utils

import (
	"strconv"
)

// Int64ToString converts an int64 to its string representation.
func Int64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}

// StringToInt64 converts a string to an int64.  Returns an error if the conversion fails.
func StringToInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}
