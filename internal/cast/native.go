package cast

import (
	"strconv"
)

func Int64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}
