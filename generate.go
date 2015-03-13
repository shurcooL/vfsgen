// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0/

package bindata

import (
	"compress/gzip"
	"fmt"
	"io"

	"github.com/shurcooL/go/vfs/httpfs/vfsutil"
)

// writeAssets writes the code file.
func writeAssets(w io.Writer, c *Config, toc []Asset) error {
	err := writeHeader(w, c)
	if err != nil {
		return err
	}

	for i := range toc {
		err = writeAsset(w, c, &toc[i])
		if err != nil {
			return err
		}
	}

	return nil
}

// writeHeader writes output file headers.
func writeHeader(w io.Writer, c *Config) error {
	_, err := fmt.Fprint(w, `import (
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"
	"unsafe"
)

func bindata_read(data, name string) ([]byte, error) {
	var empty [0]byte
	sx := (*reflect.StringHeader)(unsafe.Pointer(&data))
	b := empty[:]
	bx := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	bx.Data = sx.Data
	bx.Len = len(data)
	bx.Cap = bx.Len

	gz, err := gzip.NewReader(bytes.NewBuffer(b))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	return buf.Bytes(), nil
}

func bindata_read_compressed(data, name string) ([]byte, error) {
	var empty [0]byte
	sx := (*reflect.StringHeader)(unsafe.Pointer(&data))
	b := empty[:]
	bx := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	bx.Data = sx.Data
	bx.Len = len(data)
	bx.Cap = bx.Len
	return b, nil
}

type asset struct {
	bytes           []byte
	compressedBytes []byte
	info            bindata_file_info
}

func (_ *asset) Close() error { return nil }

type bindata_file_info struct {
	name             string
	uncompressedSize int64
	compressedSize   int64
	mode             os.FileMode
	modTime          time.Time
}

func (fi bindata_file_info) Name() string {
	return fi.name
}
func (fi bindata_file_info) Mode() os.FileMode {
	return fi.mode
}
func (fi bindata_file_info) ModTime() time.Time {
	return fi.modTime
}
func (fi bindata_file_info) IsDir() bool {
	return false
}
func (fi bindata_file_info) Sys() interface{} {
	return nil
}

type uncompressedFileInfo struct{ bindata_file_info }

func (fi uncompressedFileInfo) Size() int64 {
	return fi.uncompressedSize
}

type compressedFileInfo struct{ bindata_file_info }

func (fi compressedFileInfo) Size() int64 {
	return fi.compressedSize
}

`)
	return err
}

// writeAsset write a entry for the given asset.
// An entry is a function which embeds and returns
// the file's byte content.
func writeAsset(w io.Writer, c *Config, asset *Asset) error {
	fd, err := c.Input.Open(asset.Path)
	if err != nil {
		return err
	}
	defer fd.Close()

	compressedSize, err := writeCompressedAsset(w, asset, fd)
	if err != nil {
		return err
	}
	return writeAssetCommon(w, c, asset, compressedSize)
}

func writeCompressedAsset(w io.Writer, asset *Asset, r io.Reader) (int64, error) {
	_, err := fmt.Fprintf(w, `var _%s = "`, asset.Func)
	if err != nil {
		return 0, err
	}

	sw := &StringWriter{Writer: w}
	gz := gzip.NewWriter(sw)
	_, err = io.Copy(gz, r)
	gz.Close()

	if err != nil {
		return 0, err
	}

	_, err = fmt.Fprintf(w, `"

func %s_bytes() ([]byte, error) {
	return bindata_read(
		_%s,
		%q,
	)
}

func %s_bytes_compressed() ([]byte, error) {
	return bindata_read_compressed(
		_%s,
		%q,
	)
}

`, asset.Func, asset.Func, asset.Name, asset.Func, asset.Func, asset.Name)
	return sw.c, err
}

func writeAssetCommon(w io.Writer, c *Config, asset *Asset, compressedSize int64) error {
	fi, err := vfsutil.Stat(c.Input, asset.Path)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, `func %s() (*asset, error) {
	bytes, err := %s_bytes()
	if err != nil {
		return nil, err
	}

	compressedBytes, err := %s_bytes_compressed()
	if err != nil {
		return nil, err
	}

	info := bindata_file_info{name: %q, uncompressedSize: %d, compressedSize: %d, mode: os.FileMode(%d), modTime: time.Unix(%d, 0)}
	a := &asset{bytes: bytes, compressedBytes: compressedBytes, info: info}
	return a, nil
}

`, asset.Func, asset.Func, asset.Func, asset.Name, fi.Size(), compressedSize, uint32(fi.Mode()), fi.ModTime().Unix())
	return err
}

func writeVFS(w io.Writer) error {
	_, err := fmt.Fprint(w, `
var fileTimestamp = time.Now()

// FakeFile implements os.FileInfo interface for a given path and size
type FakeFile struct {
	// Path is the path of this file
	Path string
	// Dir marks of the path is a directory
	Dir bool
	// Len is the length of the fake file, zero if it is a directory
	Len int64
}

func (f *FakeFile) Name() string {
	_, name := filepath.Split(f.Path)
	return name
}

func (f *FakeFile) Mode() os.FileMode {
	mode := os.FileMode(0644)
	if f.Dir {
		return mode | os.ModeDir
	}
	return mode
}

func (f *FakeFile) ModTime() time.Time {
	return fileTimestamp
}

func (f *FakeFile) Size() int64 {
	return f.Len
}

func (f *FakeFile) IsDir() bool {
	return f.Mode().IsDir()
}

func (f *FakeFile) Sys() interface{} {
	return nil
}

// AssetFile implements http.File interface for a no-directory file with content
type AssetFile struct {
	*bytes.Reader
	*asset
}

/*func NewAssetFile(name string, content []byte) *AssetFile {
	return &AssetFile{
		bytes.NewReader(content),
		FakeFile{name, false, int64(len(content))},
	}
}*/

func (f *AssetFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, errors.New("not a directory")
}

func (f *AssetFile) Stat() (os.FileInfo, error) {
	return uncompressedFileInfo{f.asset.info}, nil
}

func (f *AssetFile) GzipBytes() []byte {
	log.Println("using GzipBytes for", f.asset.info.Name())
	return f.asset.compressedBytes
}

func (_ *AssetFile) Close() error { return nil }

type AssetFileOld struct {
	*bytes.Reader
	FakeFile
}

func (f *AssetFileOld) Readdir(count int) ([]os.FileInfo, error) {
	return nil, errors.New("not a directory")
}

func (f *AssetFileOld) Stat() (os.FileInfo, error) {
	return f, nil
}

func (_ *AssetFileOld) Close() error {
	return nil
}

// AssetDirectory implements http.File interface for a directory
type AssetDirectory struct {
	name          string
	io.ReadSeeker // TODO: nil so will panic.
	childrenRead  int
	children      []os.FileInfo
}

func NewAssetDirectory(name string, children []string, fs *AssetFS) *AssetDirectory {
	fileinfos := make([]os.FileInfo, 0, len(children))
	for _, child := range children {
		_, err := AssetDir(filepath.Join(name, child))
		fileinfos = append(fileinfos, &FakeFile{child, err == nil, 0})
	}
	return &AssetDirectory{
		/*AssetFileOld: AssetFileOld{
			bytes.NewReader(nil),
			FakeFile{Path: name, Dir: true, Len:0},
		},*/
		name:         name,
		childrenRead: 0,
		children:     fileinfos,
	}
}

func (f *AssetDirectory) Readdir(count int) ([]os.FileInfo, error) {
	if count <= 0 {
		return f.children, nil
	}
	if f.childrenRead+count > len(f.children) {
		count = len(f.children) - f.childrenRead
	}
	rv := f.children[f.childrenRead : f.childrenRead+count]
	f.childrenRead += count
	return rv, nil
}

func (f *AssetDirectory) Stat() (os.FileInfo, error) {
	return &FakeFile{Path: f.name, Dir: true, Len: 0}, nil
}

func (_ *AssetDirectory) Close() error { return nil }

// TODO: To be final output.
//var AssetsFs = godocfs.New(&AssetFS{})
var AssetsFs http.FileSystem = &AssetFS{}

// AssetFS implements http.FileSystem, allowing
// embedded files to be served from net/http package.
type AssetFS struct{}

func (fs *AssetFS) Open(name string) (http.File, error) {
	if len(name) > 0 && name[0] == '/' {
		name = name[1:]
	}
	if children, err := AssetDir(name); err == nil {
		return NewAssetDirectory(name, children, fs), nil
	}
	a, err := Asset2(name)
	if err != nil {
		return nil, err
	}
	//return a, nil
	return &AssetFile{
		Reader: bytes.NewReader(a.bytes),
		asset:  a,
		//FakeFile: FakeFile{name, false, int64(len(a.bytes))},
	}, nil
}
`)
	return err
}
