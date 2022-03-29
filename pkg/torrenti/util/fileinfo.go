package util

import (
	"bytes"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
)

var (
	_ fs.FileInfo = &File{}
	_ fs.File     = &File{}
)

type File struct {
	Path     string
	Length   int64
	FileMode os.FileMode
	Modified time.Time
	Internal any
	Data     []byte

	URL      string // file origin url
	Response *http.Response
	Reader   io.ReadCloser
}

func (f *File) ArchiveEntry(file *File) *File {
	file.URL = f.URL
	return file
}

func (f *File) Stat() (fs.FileInfo, error) {
	return f, nil
}

func (f *File) Read(i []byte) (int, error) {
	if f.Reader != nil {
		return f.Reader.Read(i)
	}
	return 0, io.EOF
}

func (f *File) Close() error {
	if f.Reader != nil {
		defer func() {
			f.Reader = nil
		}()
		return f.Reader.Close()
	}
	return nil
}

func (f File) ReadAll() (data []byte, err error) {
	if f.Data != nil {
		return f.Data, nil
	}
	if _, err = os.Stat(f.Path); err == nil {
		f.Data, err = os.ReadFile(f.Path)
		err = errors.Wrapf(err, "os.ReadFile: %v", f.Path)
	} else if f.URL != "" {
		f.Response, err = http.Get(f.URL)
		err = errors.Wrapf(err, "http.Get: %v", f.URL)
		if err == nil {
			defer f.Response.Body.Close()
			f.Data, err = io.ReadAll(f.Response.Body)
			err = errors.Wrapf(err, "io.ReadAll http response: %v", f.URL)
		}
	} else {
		err = errors.New("unable to load file")
	}
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
