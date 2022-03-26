package handlers

import (
	"bytes"
	"context"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/wenerme/torrenti/pkg/torrenti/util"
	"golang.org/x/exp/slices"
)

type Handler struct {
	Name       string
	Extensions []string
	Handle     func(ctx context.Context, in *util.File, cb func(ctx context.Context, in *util.File) error) error
}

func Ext(f *util.File) string {
	ext := strings.ToLower(filepath.Ext(f.Name()))
	neo := ext
	switch ext {
	case ".rar", ".zip":
		// common mistake
		dt := http.DetectContentType(f.Data)
		switch dt {
		case "application/x-rar-compressed":
			neo = ".rar"
		case "application/zip":
			neo = ".zip"
		}
	default:
		switch {
		case ext != ".torrent" && isTorrent(f.Data):
			neo = ".torrent"
		}
	}
	if neo != ext {
		log.Debug().Str("file", f.Path).Str("ext", neo).Msg("fix ext")
		ext = neo
	}
	return ext
}

func isTorrent(d []byte) bool {
	return bytes.HasPrefix(d, []byte("d8:announce")) || bytes.HasPrefix(d, []byte("d13:announce-list"))
}

var subtitleExts = []string{".srt", ".ass", ".vtt", ".ssa", ".ttml", ".tml", ".sami", ".sbv", ".sub", ".rt", ".scc", ".dfxp", ".sup", ".idx", ".mks"}

func init() {
	slices.Sort(subtitleExts)
}

func IsSubtitleExt(ext string) bool {
	return slices.BinarySearch(subtitleExts, ext) >= 0
}
