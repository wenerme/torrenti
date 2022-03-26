package util

import (
	"strings"
	"time"
)

type BuildInfo struct {
	Modified  bool
	Time      time.Time
	Revision  string
	VCS       string
	CGO       bool
	GOOS      string
	GOARCH    string
	GoVersion string
	Version   string
}

func (v BuildInfo) String() string {
	sb := &strings.Builder{}
	var tags []string

	sb.WriteString(v.Version)

	if v.VCS != "" {
		sb.WriteString(" ")
		sb.WriteString(v.VCS)
		sb.WriteString(" ")
		sb.WriteString(v.Revision[:7])
		sb.WriteString(" ")
		sb.WriteString(v.Time.Format("2006-01-02 15:04:05"))
	} else {
		tags = append(tags, "novcs")
	}
	if v.Modified {
		tags = append(tags, "dirty")
	}
	if v.CGO {
		tags = append(tags, "cgo")
	}
	tags = append(tags, v.GoVersion, v.GOOS, v.GOARCH)
	if len(tags) > 0 {
		sb.WriteString(" (")
		sb.WriteString(strings.Join(tags, ", "))
		sb.WriteString(")")
	}
	return sb.String()
}
