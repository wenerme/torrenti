package nilx

import (
	"strings"
)

func EmptyToNil(s ...string) *string {
	for _, v := range s {
		if v != "" {
			return &v
		}
	}
	return nil
}

func TrimEmptyToNil(s string) *string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	return &s
}

func NilToEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func ZeroToNil(v int) *int {
	if v == 0 {
		return nil
	}
	return &v
}
