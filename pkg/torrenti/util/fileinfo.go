package util

import (
	"bytes"
	"io"
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
	Internal any
	Data     []byte
}

func (f File) ReadAll() ([]byte, error) {
	return f.Data, nil
}

func (f File) Open() (io.ReadCloser, error) {
	if f.Data != nil {
		return io.NopCloser(bytes.NewReader(f.Data)), nil
	}
	return os.Open(f.Path)
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
	return f.FileMode.IsDir()
}

func (f File) Sys() any {
	return f.Internal
}
