package protox

import (
	"time"

	"google.golang.org/protobuf/types/known/durationpb"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func UnixToTimestamp(t int64) *timestamppb.Timestamp {
	if t == 0 {
		return nil
	}
	return &timestamppb.Timestamp{
		Seconds: t,
	}
}

func UnixMilliToTimestamp(t int64) *timestamppb.Timestamp {
	if t == 0 {
		return nil
	}
	return ToTimestamp(time.UnixMilli(t))
}

func ToTimestamp(t time.Time) *timestamppb.Timestamp {
	if t.IsZero() {
		return nil
	}
	return timestamppb.New(t)
}

func PtrToTimestamp(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(*t)
}

var (
	FromTimestamp = timestamppb.New
	FromDuration  = durationpb.New
)
