package archives

import (
	"bytes"
	"context"
	"io"

	"github.com/nwaples/rardecode/v2"
	"github.com/wenerme/torrenti/pkg/torrenti/util"
)

func Unrar(ctx context.Context, in *util.File, cb func(context.Context, *util.File) error) (err error) {
	r, err := rardecode.NewReader(bytes.NewReader(in.Data))
	if err != nil {
		return err
	}

	var next *rardecode.FileHeader
	for {
		next, err = r.Next()
		if err == io.EOF {
			err = nil
			return nil
		}
		if err != nil {
			return err
		}

		f := in.ArchiveEntry(&util.File{
			Path:     next.Name,
			Length:   next.UnPackedSize,
			FileMode: next.Mode(),
			Modified: next.ModificationTime,
			Internal: next,
		})
		f.Data, err = io.ReadAll(r)
		if err != nil {
			return err
		}
		err = cb(ctx, f)
		if err != nil {
			return err
		}
	}
}
