package timeprimitive

import "time"

type Timestamp struct {
	value time.Time
}

func Now() Timestamp {
	return Timestamp{value: normalize(time.Now())}
}

func TimestampFrom(value time.Time) Timestamp {
	return Timestamp{value: normalize(value)}
}

func (t Timestamp) Time() time.Time {
	return t.value
}

func (t Timestamp) IsZero() bool {
	return t.value.IsZero()
}

func (t Timestamp) String() string {
	return t.value.Format(time.RFC3339Nano)
}

func normalize(value time.Time) time.Time {
	return value.UTC().Truncate(time.Microsecond)
}
