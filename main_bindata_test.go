// generated via `go generate`; do not edit

package bindata_test

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"
	"unsafe"
)

// For assetfs.
import (
	"errors"
	"net/http"
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

var _not_worth_compressing_file_txt = "\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xf2\x2c\x29\x56\xc8\xcb\x2f\xca\x4d\xcc\x51\x48\xce\xcf\x2b\x49\xcd\x03\xf2\x13\x8b\x52\x15\x32\x52\x8b\x52\xf5\x00\x01\x00\x00\xff\xff\xdc\xc7\xff\x13\x1d\x00\x00\x00"

func not_worth_compressing_file_txt_bytes() ([]byte, error) {
	return bindata_read(
		_not_worth_compressing_file_txt,
		"not-worth-compressing-file.txt",
	)
}

func not_worth_compressing_file_txt_bytes_compressed() ([]byte, error) {
	return bindata_read_compressed(
		_not_worth_compressing_file_txt,
		"not-worth-compressing-file.txt",
	)
}

func not_worth_compressing_file_txt() (*asset, error) {
	bytes, err := not_worth_compressing_file_txt_bytes()
	if err != nil {
		return nil, err
	}

	compressedBytes, err := not_worth_compressing_file_txt_bytes_compressed()
	if err != nil {
		return nil, err
	}

	info := bindata_file_info{name: "not-worth-compressing-file.txt", uncompressedSize: 29, compressedSize: 52, mode: os.FileMode(292), modTime: time.Unix(-62135596800, 0)}
	a := &asset{bytes: bytes, compressedBytes: compressedBytes, info: info}
	return a, nil
}

var _sample_file_txt = "\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x0a\xc9\xc8\x2c\x56\x48\xcb\xcc\x49\x55\x48\xce\xcf\x2d\x28\x4a\x2d\x2e\x4e\x2d\x56\x28\x4f\xcd\xc9\xd1\x53\x70\xca\x49\x1c\xd4\x20\x43\x0f\x10\x00\x00\xff\xff\x76\x5a\x3e\xaa\xbd\x00\x00\x00"

func sample_file_txt_bytes() ([]byte, error) {
	return bindata_read(
		_sample_file_txt,
		"sample-file.txt",
	)
}

func sample_file_txt_bytes_compressed() ([]byte, error) {
	return bindata_read_compressed(
		_sample_file_txt,
		"sample-file.txt",
	)
}

func sample_file_txt() (*asset, error) {
	bytes, err := sample_file_txt_bytes()
	if err != nil {
		return nil, err
	}

	compressedBytes, err := sample_file_txt_bytes_compressed()
	if err != nil {
		return nil, err
	}

	info := bindata_file_info{name: "sample-file.txt", uncompressedSize: 189, compressedSize: 58, mode: os.FileMode(292), modTime: time.Unix(-62135596800, 0)}
	a := &asset{bytes: bytes, compressedBytes: compressedBytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

func Asset2(name string) (*asset, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

func AssetCompressed(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.compressedBytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return uncompressedFileInfo{a.info}, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"not-worth-compressing-file.txt": not_worth_compressing_file_txt,
	"sample-file.txt": sample_file_txt,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.children))
	for name := range node.children {
		rv = append(rv, name)
	}
	return rv, nil
}

type _bintree_t struct {
	Func     func() (*asset, error)
	children map[string]*_bintree_t
}

var _bintree = &_bintree_t{nil, map[string]*_bintree_t{
	"not-worth-compressing-file.txt": &_bintree_t{not_worth_compressing_file_txt, map[string]*_bintree_t{
	}},
	"sample-file.txt": &_bintree_t{sample_file_txt, map[string]*_bintree_t{
	}},
}}

var (
	fileTimestamp = time.Now()
)

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
	log.Println("using GzipBytes!")
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
