package handlers

import (
	"bytes"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/wenerme/torrenti/pkg/torrenti/util"
)

func Ext(f *util.File) string {
	ext := strings.ToLower(filepath.Ext(f.Name()))
	neo := ext
	// sanitize
	switch ext {
	case ".rar", ".zip":
		// common mistake
		if f.Data == nil {
			break
		}
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
		case len(ext) == 1 || ext[0] != '.':
			ext = ""
		// case strings.Contains(".abcdefghijklmnopqrstuvwxyz", ext):
		case !isValidExt(ext):
			ext = ""
		case ext != ".html" && ext != ".htm" && len(f.Data) > 4 && strings.EqualFold(string(f.Data[:5]), "<html"):
			// error page
			ext = ".html"
		}
	}
	if neo != ext {
		log.Debug().Str("file", f.Path).Str("ext", neo).Msg("fix ext")
		ext = neo
	}
	return ext
}

func isValidExt(s string) bool {
	for _, r := range s {
		switch {
		case r == '.':
		case r >= 'a' && r <= 'z':
		default:
			return false
		}
	}
	return true
}

func isTorrent(d []byte) bool {
	return bytes.HasPrefix(d, []byte("d8:announce")) || bytes.HasPrefix(d, []byte("d13:announce-list"))
}
