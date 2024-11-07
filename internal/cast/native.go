package cast

import (
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/protobuf/types/known/durationpb"
)

func Int64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}

// Convert pgtype.Interval to *durationpb.Duration.
func PgIntervalToDuration(interval pgtype.Interval) *durationpb.Duration {
	return durationpb.New(time.Duration(interval.Microseconds) * time.Microsecond)
}

// Convert *durationpb.Duration to pgtype.Interval.
func DurationToPgInterval(dur *durationpb.Duration) pgtype.Interval {
	return pgtype.Interval{
		Microseconds: dur.AsDuration().Microseconds(),
		Valid:        true,
	}
}
