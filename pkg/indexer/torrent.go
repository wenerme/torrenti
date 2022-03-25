package indexer

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/wenerme/torrenti/pkg/magnet"
	"github.com/xgfone/bt/metainfo"
)

func ParseTorrent(v string) (o *Torrent, err error) {
	o = &Torrent{}
	fi, _ := os.Stat(v)

	switch {
	case fi != nil:
		o.FileInfo = fi
		o.File, err = filepath.Abs(v)
	case strings.HasPrefix(v, "http:") || strings.HasPrefix(v, "https:"):
		o.URL = v
	case strings.HasPrefix(v, "magnet:"):
		o.Magnet, err = magnet.Parse(v)
	case !strings.ContainsAny(v, "./\\"): // hash
		o.Hash, err = magnet.ParseHash(v)
		o.Magnet = magnet.Magnet{Hash: o.Hash}
	default:
		err = fmt.Errorf("invalid torrent: %q", v)
	}
	return
}

type Torrent struct {
	Magnet   magnet.Magnet
	Hash     magnet.Hash
	File     string
	FileInfo fs.FileInfo
	Data     []byte
	Meta     *metainfo.MetaInfo
	URL      string
}

func (t *Torrent) Load() (err error) {
	switch {
	case t.URL != "":
		var resp *http.Response
		resp, err = http.Get(t.URL)
		if err == nil {
			defer resp.Body.Close()
			t.Data, err = io.ReadAll(resp.Body)
		}
		if err != nil {
			return
		}

		_, params, _ := mime.ParseMediaType(resp.Header.Get("Content-Disposition"))
		fi := &fileInfo{}
		fi.name = params["filename"]
		fi.size = int64(len(t.Data))
		t.FileInfo = fi

	case t.File != "":
		var f *os.File

		f, err = os.Open(t.File)
		if err == nil {
			defer f.Close()
			t.Data, err = io.ReadAll(f)

		}
		if err != nil {
			return
		}

	case !t.Magnet.Hash.IsZero():
		err = errors.New("TODO: load magnet")
	default:
		err = errors.New("invalid torrent info")
	}
	if err == nil {
		err = t.loadData()
	}
	return
}

func (t *Torrent) loadData() (err error) {
	var meta metainfo.MetaInfo

	meta, err = metainfo.Load(bytes.NewReader(t.Data))
	t.Meta = &meta
	t.Hash = magnet.Hash{
		Digest: meta.InfoHash().Bytes(),
	}
	t.Magnet = magnet.Magnet{
		Hash: t.Hash,
	}
	return
}

var _ fs.FileInfo = fileInfo{}

type fileInfo struct {
	name string
	size int64
	mode os.FileMode
	mod  time.Time
	dir  bool
	sys  any
}

func (f fileInfo) Name() string {
	return f.name
}

func (f fileInfo) Size() int64 {
	return f.size
}

func (f fileInfo) Mode() fs.FileMode {
	return f.mode
}

func (f fileInfo) ModTime() time.Time {
	return f.mod
}

func (f fileInfo) IsDir() bool {
	return f.dir
}

func (f fileInfo) Sys() any {
	return f.sys
}
