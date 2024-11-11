package db

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/protobuf/types/known/durationpb"
)

// IntervalToDuration converts a pgtype.Interval to a time.Duration.
func IntervalToDuration(interval pgtype.Interval) time.Duration {
	ms := interval.Microseconds +
		int64(interval.Days*24*60*60*1000*1000) +
		int64(interval.Months*30*24*60*60*1000*1000)

	return time.Duration(ms) * time.Microsecond
}

// IntervalToProto converts a pgtype.Interval to a durationpb.Duration.
func IntervalToProto(interval pgtype.Interval) *durationpb.Duration {
	return durationpb.New(IntervalToDuration(interval))
}

// DurationToInterval converts a time.Duration to a pgtype.Interval.
func DurationToInterval(d time.Duration) pgtype.Interval {
	return pgtype.Interval{
		Microseconds: int64(d / time.Microsecond),
	}
}

// ProtoToInterval converts a durationpb.Duration to a pgtype.Interval.
func ProtoToInterval(d *durationpb.Duration) pgtype.Interval {
	return pgtype.Interval{
		Microseconds: int64(d.AsDuration() / time.Microsecond),
	}
}
