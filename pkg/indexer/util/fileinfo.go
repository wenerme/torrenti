package util

import (
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

var _ fs.FileInfo = File{}

type File struct {
	Path     string
	Length   int64
	FileMode os.FileMode
	Modified time.Time
	Dir      bool
	Internal any
	Data     []byte
}

func (f File) Name() string {
	return filepath.Base(f.Path)
}

func (f File) Size() int64 {
	return f.Length
}

func (f File) Mode() fs.FileMode {
	return f.FileMode
}

func (f File) ModTime() time.Time {
	return f.Modified
}

func (f File) IsDir() bool {
	return f.Dir
}

func (f File) Sys() any {
	return f.Internal
}
