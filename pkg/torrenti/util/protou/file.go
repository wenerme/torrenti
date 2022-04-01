package protou

import (
	"os"

	"github.com/wenerme/torrenti/pkg/apis/media/common"
	"github.com/wenerme/torrenti/pkg/torrenti/util"
)

func ToFile(f *common.File) *util.File {
	out := &util.File{
		Path:     f.GetPath(),
		Length:   f.GetLength(),
		FileMode: os.FileMode(f.GetFileMode()),
		Modified: f.Modified.AsTime(),
		Internal: f,
		Data:     f.Data,
		URL:      f.GetUrl(),
		Response: nil,
	}
	return out
}
