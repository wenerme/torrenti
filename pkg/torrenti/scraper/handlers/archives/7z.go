package archives

import (
	"bytes"
	"context"
	"io"

	"github.com/bodgit/sevenzip"
	"github.com/wenerme/torrenti/pkg/torrenti/util"
)

func Un7z(ctx context.Context, in *util.File, cb func(context.Context, *util.File) error) error {
	r, err := sevenzip.NewReader(bytes.NewReader(in.Data), in.Length)
	if err != nil {
		return err
	}
	for _, fe := range r.File {
		if fe.Mode().IsDir() {
			continue
		}
		f := &util.File{
			Path:     fe.Name,
			FileMode: fe.Mode(),
			Modified: fe.Modified,
			Length:   int64(fe.UncompressedSize),
			Internal: fe,
		}
		var rr io.ReadCloser
		rr, err = fe.Open()
		if err == nil {
			f.Data, err = io.ReadAll(rr)
			if err == nil {
				err = rr.Close()
			}
		}
		if err != nil {
			return err
		}
		if err = cb(ctx, f); err != nil {
			return err
		}
	}
	return err
}
