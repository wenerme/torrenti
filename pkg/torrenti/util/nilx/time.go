package nilx

import "time"

func TimeNilToZero(s *time.Time) time.Time {
	if s == nil {
		return time.Time{}
	}
	return *s
}

func TimeZeroToNil(s time.Time) *time.Time {
	if s.IsZero() {
		return nil
	}
	return &s
}
