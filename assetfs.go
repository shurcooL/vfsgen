package bindata

import (
	"fmt"
	"io"
)

func writeVFS(w io.Writer) error {
	_, err := fmt.Fprintf(w, `
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
