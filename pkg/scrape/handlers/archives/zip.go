package archives

import (
	"archive/zip"
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/wenerme/torrenti/pkg/torrenti/util"
	"golang.org/x/text/encoding/simplifiedchinese"
)

func Unzip(ctx context.Context, in *util.File, cb func(context.Context, *util.File) error) (err error) {
	var zr *zip.Reader
	var data []byte
	data, err = in.ReadAll()
	if err != nil {
		return
	}

	zr, err = zip.NewReader(bytes.NewReader(data), in.Length)
	if err != nil {
		fn := filepath.Join(os.TempDir(), in.Name())
		_ = os.WriteFile(fn, data, 0o644)
		log.Warn().Str("path", fn).Msg("file dump")
		return
	}
	// gbk := simplifiedchinese.GB18030.NewDecoder()
	gbk := simplifiedchinese.GBK.NewDecoder()

	for _, fe := range zr.File {
		fi := &util.File{
			Path:     fe.Name,
			FileMode: fe.Mode(),
			Modified: fe.Modified,
			Length:   int64(fe.UncompressedSize64),
			Internal: fe,
		}
		if fe.NonUTF8 {
			// gbk
			fi.Path, err = gbk.String(fe.Name)
			if err != nil {
				return errors.Wrap(err, "decode gbk")
			}
			// var dr *chardet.Result
			// dr, err = chardet.NewTextDetector().DetectBest([]byte(filepath.Base(fe.Name)))
			// log.Debug().Str("zip", in.Path).Str("path", fi.Path).Str("charset", dr.Charset).Msg("unzip gbk decode")
			log.Debug().Str("zip", in.Path).Str("path", fi.Path).Msg("unzip gbk decode")
		}
		if fi.FileMode.IsDir() {
			continue
		}
		var ff io.ReadCloser
		ff, err = fe.Open()
		if err != nil {
			return errors.Wrapf(err, "zip open entry: utf8=%v %q", !fe.NonUTF8, fi.Path)
		}
		fi.Data, err = io.ReadAll(ff)
		ff.Close()
		if err != nil {
			return err
		}
		err = cb(ctx, fi)
		if err != nil {
			return err
		}
	}
	return
}
