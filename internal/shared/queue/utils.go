package queue

import (
	"strings"
)

func format(args ...string) string {
	return strings.Join(args, ".")
}
