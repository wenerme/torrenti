package torrenti

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"github.com/wenerme/torrenti/pkg/magnet"
	"github.com/wenerme/torrenti/pkg/torrenti/util"
	"github.com/xgfone/bt/metainfo"
)

func ParseTorrent(v string) (o *Torrent, err error) {
	o = &Torrent{}
	fi, _ := os.Stat(v)

	switch {
	case fi != nil:
		o.FileInfo = fi
		o.Path, err = filepath.Abs(v)
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
	Path     string
	FileInfo fs.FileInfo
	Data     []byte
	Meta     *metainfo.MetaInfo
	URL      string
	Response *http.Response
	File     *util.File
	H        string
}

func (t *Torrent) Load() (err error) {
	switch {
	case t.URL != "":
		if t.Response == nil && t.Data == nil {
			t.Response, err = http.Get(t.URL)
			if err != nil {
				return
			}
		}
		resp := t.Response
		if t.Data == nil {
			defer resp.Body.Close()
			t.Data, err = io.ReadAll(resp.Body)
		}
		if err != nil {
			return
		}
		if t.FileInfo == nil {
			_, params, _ := mime.ParseMediaType(resp.Header.Get("Content-Disposition"))
			fi := &util.File{}
			fi.Path = params["filename"]
			fi.Length = int64(len(t.Data))
			t.FileInfo = fi
		}
	case t.Path != "":
		var f *os.File

		f, err = os.Open(t.Path)
		if err == nil {
			defer f.Close()
			t.Data, err = io.ReadAll(f)

		}
		if err != nil {
			return
		}

	case !t.Magnet.Hash.IsZero():
		// dht.Config{}
		err = errors.New("TODO: load magnet")
	case t.Data != nil && t.FileInfo != nil:
		// data ready
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
	err = errors.Wrap(err, "load metainfo")
	t.Meta = &meta
	t.Hash = magnet.Hash{
		Digest: meta.InfoHash().Bytes(),
	}
	t.Magnet = magnet.Magnet{
		Hash: t.Hash,
	}
	return
}
