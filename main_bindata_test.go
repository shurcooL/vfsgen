// generated via `go generate`; do not edit

package bindata_test

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

var AssetsFs http.FileSystem = _assetFS

type AssetFS map[string]interface{}

var _assetFS = AssetFS{
	"/folderA/file1.txt": &compressedFile{
		name:              "file1.txt",
		compressedContent: []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x0a\x2e\x29\x4d\x4b\xd3\x03\x04\x00\x00\xff\xff\x27\xbb\x40\xc8\x06\x00\x00\x00"),
		uncompressedSize:  6,
		modTime:           time.Time{},
	},
	"/folderA/file2.txt": &compressedFile{
		name:              "file2.txt",
		compressedContent: []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x0a\x2e\x29\x4d\x4b\xd3\x03\x04\x00\x00\xff\xff\x27\xbb\x40\xc8\x06\x00\x00\x00"),
		uncompressedSize:  6,
		modTime:           time.Time{},
	},
	"/folderB/folderC/file3.txt": &compressedFile{
		name:              "file3.txt",
		compressedContent: []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x0a\x2e\x29\x4d\x4b\xd3\x03\x04\x00\x00\xff\xff\x27\xbb\x40\xc8\x06\x00\x00\x00"),
		uncompressedSize:  6,
		modTime:           time.Time{},
	},
	"/not-worth-compressing-file.txt": &compressedFile{
		name:              "not-worth-compressing-file.txt",
		compressedContent: []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xf2\x2c\x29\x56\xc8\xcb\x2f\xca\x4d\xcc\x51\x48\xce\xcf\x2b\x49\xcd\x03\xf2\x13\x8b\x52\x15\x32\x52\x8b\x52\xf5\x00\x01\x00\x00\xff\xff\xdc\xc7\xff\x13\x1d\x00\x00\x00"),
		uncompressedSize:  29,
		modTime:           time.Time{},
	},
	"/sample-file.txt": &compressedFile{
		name:              "sample-file.txt",
		compressedContent: []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x0a\xc9\xc8\x2c\x56\x48\xcb\xcc\x49\x55\x48\xce\xcf\x2d\x28\x4a\x2d\x2e\x4e\x2d\x56\x28\x4f\xcd\xc9\xd1\x53\x70\xca\x49\x1c\xd4\x20\x43\x0f\x10\x00\x00\xff\xff\x76\x5a\x3e\xaa\xbd\x00\x00\x00"),
		uncompressedSize:  189,
		modTime:           time.Time{},
	},
}

func init() {
	_assetFS["/folderB/folderC"] = &dir{
		name: "folderC",
		entries: []os.FileInfo{
			_assetFS["/folderB/folderC/file3.txt"].(os.FileInfo),
		},
		modTime: time.Time{},
	}
	_assetFS["/folderA"] = &dir{
		name: "folderA",
		entries: []os.FileInfo{
			_assetFS["/folderA/file1.txt"].(os.FileInfo),
			_assetFS["/folderA/file2.txt"].(os.FileInfo),
		},
		modTime: time.Time{},
	}
	_assetFS["/folderB"] = &dir{
		name: "folderB",
		entries: []os.FileInfo{
			_assetFS["/folderB/folderC"].(os.FileInfo),
		},
		modTime: time.Time{},
	}
	_assetFS["/"] = &dir{
		name: "/",
		entries: []os.FileInfo{
			_assetFS["/folderA"].(os.FileInfo),
			_assetFS["/folderB"].(os.FileInfo),
			_assetFS["/not-worth-compressing-file.txt"].(os.FileInfo),
			_assetFS["/sample-file.txt"].(os.FileInfo),
		},
		modTime: time.Time{},
	}
}

func (fs AssetFS) Open(path string) (http.File, error) {
	f, ok := fs[path]
	if !ok {
		return nil, os.ErrNotExist
	}

	if cf, ok := f.(*compressedFile); ok {
		gr, err := gzip.NewReader(bytes.NewReader(cf.compressedContent))
		if err != nil {
			// This should never happen because we generate the gzip bytes such that they are always valid.
			panic("unexpected error reading own gzip compressed bytes: " + err.Error())
		}
		return &compressedFileInstance{
			compressedFile: cf,
			gr:             gr,
		}, nil
	}

	return f.(http.File), nil
}

// compressedFile is ...
type compressedFile struct {
	name              string
	compressedContent []byte
	uncompressedSize  int64
	modTime           time.Time
}

func (f *compressedFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, fmt.Errorf("cannot Readdir from file %s", f.name)
}
func (f *compressedFile) Stat() (os.FileInfo, error) { return f, nil }

func (f *compressedFile) GzipBytes() []byte {
	log.Println("using GzipBytes for", f.name)
	return f.compressedContent
}

func (f *compressedFile) Name() string       { return f.name }
func (f *compressedFile) Size() int64        { return f.uncompressedSize }
func (f *compressedFile) Mode() os.FileMode  { return 0444 }
func (f *compressedFile) ModTime() time.Time { return f.modTime }
func (f *compressedFile) IsDir() bool        { return false }
func (f *compressedFile) Sys() interface{}   { return nil }

type compressedFileInstance struct {
	*compressedFile
	gr io.ReadCloser
}

func (f *compressedFileInstance) Read(p []byte) (n int, err error) {
	return f.gr.Read(p)
}
func (f *compressedFileInstance) Seek(offset int64, whence int) (int64, error) {
	panic("Seek not yet implemented")
}
func (f *compressedFileInstance) Close() error {
	return f.gr.Close()
}

// dir is ...
type dir struct {
	name    string
	entries []os.FileInfo
	modTime time.Time
}

func (d *dir) Read([]byte) (int, error) {
	return 0, fmt.Errorf("cannot Read from directory %s", d.name)
}
func (d *dir) Seek(offset int64, whence int) (int64, error) {
	return 0, fmt.Errorf("cannot Seek in directory %s", d.name)
}
func (d *dir) Close() error { return nil }
func (d *dir) Readdir(count int) ([]os.FileInfo, error) {
	if count != 0 {
		log.Panicln("httpDir.Readdir count unsupported value:", count)
	}
	return d.entries, nil
}
func (d *dir) Stat() (os.FileInfo, error) { return d, nil }

func (d *dir) Name() string       { return d.name }
func (d *dir) Size() int64        { return 0 }
func (d *dir) Mode() os.FileMode  { return 0755 | os.ModeDir }
func (d *dir) ModTime() time.Time { return d.modTime }
func (d *dir) IsDir() bool        { return true }
func (d *dir) Sys() interface{}   { return nil }
