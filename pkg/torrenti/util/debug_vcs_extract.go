package util

import (
	"runtime"
	"runtime/debug"
	"time"
)

func ReadBuildInfo() (o BuildInfo) {
	bi, _ := debug.ReadBuildInfo()
	if bi == nil {
		bi = &debug.BuildInfo{
			GoVersion: runtime.Version(),
			Main: debug.Module{
				Version: "dev",
			},
		}
	}
	o.Version = bi.Main.Version
	o.GoVersion = bi.GoVersion
	o.GOOS = runtime.GOOS
	o.GOARCH = runtime.GOARCH

	for _, v := range bi.Settings {
		switch v.Key {
		case "vcs.modified":
			o.Modified = v.Value == "true"
		case "vcs.time":
			o.Time, _ = time.Parse(time.RFC3339, v.Value)
		case "vcs.revision":
			o.Revision = v.Value
		case "vcs":
			o.VCS = v.Value
		case "CGO_ENABLED":
			o.CGO = v.Value == "1"
		}
	}
	return
}
