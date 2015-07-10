package vfsgen

import (
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	pathpkg "path"
	"sort"
	"text/template"

	"github.com/shurcooL/go/vfs/httpfs/vfsutil"
)

// Generate reads assets from an input directory, converts them
// to Go code and writes new files to the output specified
// in the given configuration.
func Generate(c Config) error {
	// Ensure our configuration has sane values.
	err := c.validate()
	if err != nil {
		return err
	}

	// Create output file.
	f, err := os.Create(c.Output)
	if err != nil {
		return err
	}
	defer f.Close()

	err = writeHeader(f, c)
	if err != nil {
		return err
	}

	var toc toc
	err = findAndWriteFiles(f, c.Input, &toc)
	if err != nil {
		return err
	}

	err = writeDirEntries(f, c, toc)
	if err != nil {
		return err
	}

	err = writeTrailer(f, toc)
	if err != nil {
		return err
	}

	// Trim any potential excess.
	cur, err := f.Seek(0, os.SEEK_CUR)
	if err != nil {
		return err
	}
	err = f.Truncate(cur)
	if err != nil {
		return err
	}

	return nil
}

type toc struct {
	dirs []pathDirInfo

	CompressedFile bool // There's at least one compressedFile.
	File           bool // There's at least one uncompressed file.
}

// readDirPaths reads the directory named by dirname and returns
// a sorted list of directory paths.
func readDirPaths(fs http.FileSystem, dirname string) ([]string, error) {
	fis, err := vfsutil.ReadDir(fs, dirname)
	if err != nil {
		return nil, err
	}
	paths := make([]string, len(fis))
	for i := range fis {
		paths[i] = pathpkg.Join(dirname, fis[i].Name())
	}
	sort.Strings(paths)
	return paths, nil
}

// findAndWriteFiles recursively finds all the file paths in the given directory tree.
// They are added to the given map as keys. Values will be safe function names
// for each file, which will be used when generating the output code.
func findAndWriteFiles(f *os.File, fs http.FileSystem, toc *toc) error {
	walkFn := func(path string, fi os.FileInfo, r io.ReadSeeker, err error) error {
		if err != nil {
			log.Printf("can't stat file %s: %v\n", path, err)
			return nil
		}

		switch fi.IsDir() {
		case false:
			asset := &fileInfo{
				name:             pathpkg.Base(path),
				uncompressedSize: fi.Size(),
				modTime:          fi.ModTime().UTC(),
			}

			marker, err := f.Seek(0, os.SEEK_CUR)
			if err != nil {
				return err
			}

			// Write _vfsgen_compressedFileInfo.
			err = writeCompressedFileInfo(f, path, asset, r)
			switch err {
			default:
				return err
			case nil:
				toc.CompressedFile = true
			// If compressed file is not smaller than original, revert and write original file.
			case errCompressedNotSmaller:
				_, err = r.Seek(0, os.SEEK_SET)
				if err != nil {
					return err
				}

				_, err = f.Seek(marker, os.SEEK_SET)
				if err != nil {
					return err
				}

				// Write _vfsgen_fileInfo.
				err = writeFileInfo(f, path, asset, r)
				if err != nil {
					return err
				}
				toc.File = true
			}
		case true:
			entries, err := readDirPaths(fs, path)
			if err != nil {
				return err
			}

			asset := &dirInfo{
				name:    pathpkg.Base(path),
				entries: entries,
				modTime: fi.ModTime().UTC(),
			}

			toc.dirs = append(toc.dirs, pathDirInfo{
				path:    path,
				dirInfo: asset,
			})

			// Write _vfsgen_dirInfo.
			err = writeDirInfo(f, path, asset)
			if err != nil {
				return err
			}
		}

		return nil
	}

	err := vfsutil.WalkFiles(fs, "/", walkFn)
	if err != nil {
		return err
	}

	return nil
}

type pathDirInfo struct {
	path    string
	dirInfo *dirInfo
}

var errCompressedNotSmaller = errors.New("compressed file is not smaller than original")

// writeCompressedFileInfo writes _vfsgen_compressedFileInfo.
// It returns errCompressedNotSmaller if compressed file is not smaller than original.
func writeCompressedFileInfo(w io.Writer, path string, asset *fileInfo, r io.Reader) error {
	_, err := fmt.Fprintf(w, "\t\t%q: &_vfsgen_compressedFileInfo{\n", path)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "\t\t\tname:              %q,\n", asset.name)
	modTimeBytes, err := asset.modTime.MarshalText()
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "\t\t\tmodTime:           mustUnmarshalTextTime(%q),\n", string(modTimeBytes))
	fmt.Fprintf(w, "\t\t\tcompressedContent: []byte(\"")
	sw := &stringWriter{Writer: w}
	gz := gzip.NewWriter(sw)
	_, err = io.Copy(gz, r)
	if err != nil {
		return err
	}
	_ = gz.Close()
	if sw.c >= asset.uncompressedSize {
		return errCompressedNotSmaller
	}
	fmt.Fprintf(w, "\"),\n")
	fmt.Fprintf(w, "\t\t\tuncompressedSize:  %d,\n", asset.uncompressedSize)
	fmt.Fprintf(w, "\t\t},\n")
	return nil
}

// Write _vfsgen_fileInfo.
func writeFileInfo(w io.Writer, path string, asset *fileInfo, r io.Reader) error {
	_, err := fmt.Fprintf(w, "\t\t%q: &_vfsgen_fileInfo{\n", path)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "\t\t\tname:    %q,\n", asset.name)
	modTimeBytes, err := asset.modTime.MarshalText()
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "\t\t\tmodTime: mustUnmarshalTextTime(%q),\n", string(modTimeBytes))
	fmt.Fprintf(w, "\t\t\tcontent: []byte(\"")
	sw := &stringWriter{Writer: w}
	_, err = io.Copy(sw, r)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "\"),\n")
	fmt.Fprintf(w, "\t\t},\n")
	return nil
}

// writeDirInfo writes _vfsgen_dirInfo.
func writeDirInfo(w io.Writer, path string, asset *dirInfo) error {
	_, err := fmt.Fprintf(w, "\t\t%q: &_vfsgen_dirInfo{\n", path)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "\t\t\tname:    %q,\n", asset.name)
	modTimeBytes, err := asset.modTime.MarshalText()
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "\t\t\tmodTime: mustUnmarshalTextTime(%q),\n", string(modTimeBytes))
	fmt.Fprintf(w, "\t\t},\n")
	return nil
}

func writeDirEntries(w io.Writer, c Config, toc toc) error {
	_, err := fmt.Fprintf(w, "\t}\n\n")
	if err != nil {
		return err
	}

	for _, pathDirInfo := range toc.dirs {
		if len(pathDirInfo.dirInfo.entries) > 0 {
			fmt.Fprintf(w, "\tfs[%q].(*_vfsgen_dirInfo).entries = []os.FileInfo{\n", pathDirInfo.path)
			for _, entry := range pathDirInfo.dirInfo.entries {
				fmt.Fprintf(w, "\t\tfs[%q].(os.FileInfo),\n", entry)
			}
			fmt.Fprintf(w, "\t}\n")
		}
	}

	_, err = fmt.Fprintf(w, "\n\treturn fs\n}()\n")
	if err != nil {
		return err
	}

	return nil
}

func writeHeader(w io.Writer, c Config) error {
	header := template.Must(template.New("").Delims("⦗⦗", "⦘⦘").Parse(`// Generated by vfsgen; do not edit.

⦗⦗with .Tags⦘⦘// +build ⦗⦗.⦘⦘

⦗⦗end⦘⦘package ⦗⦗.Package⦘⦘

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// ⦗⦗.OutputName⦘⦘ statically implements the virtual filesystem given to vfsgen as input.
var ⦗⦗.OutputName⦘⦘ = func() http.FileSystem {
	mustUnmarshalTextTime := func(text string) time.Time {
		var t time.Time
		err := t.UnmarshalText([]byte(text))
		if err != nil {
			panic(err)
		}
		return t
	}

	fs := _vfsgen_fs{
`))
	return header.Execute(w, c)
}

func writeTrailer(w io.Writer, toc toc) error {
	trailer := template.Must(template.New("").Delims("⦗⦗", "⦘⦘").Parse(`
type _vfsgen_fs map[string]interface{}

func (fs _vfsgen_fs) Open(path string) (http.File, error) {
	// TODO: Maybe clean path?
	f, ok := fs[path]
	if !ok {
		return nil, &os.PathError{Op: "open", Path: path, Err: os.ErrNotExist}
	}

	switch f := f.(type) {⦗⦗if .CompressedFile⦘⦘
	case *_vfsgen_compressedFileInfo:
		gr, err := gzip.NewReader(bytes.NewReader(f.compressedContent))
		if err != nil {
			// This should never happen because we generate the gzip bytes such that they are always valid.
			panic("unexpected error reading own gzip compressed bytes: " + err.Error())
		}
		return &_vfsgen_compressedFile{
			_vfsgen_compressedFileInfo: f,
			gr: gr,
		}, nil⦗⦗end⦘⦘⦗⦗if .File⦘⦘
	case *_vfsgen_fileInfo:
		return &_vfsgen_file{
			_vfsgen_fileInfo: f,
			Reader:           bytes.NewReader(f.content),
		}, nil⦗⦗end⦘⦘
	case *_vfsgen_dirInfo:
		return &_vfsgen_dir{
			_vfsgen_dirInfo: f,
		}, nil
	default:
		// This should never happen because we generate only the above types.
		panic(fmt.Sprintf("unexpected type %T", f))
	}
}
⦗⦗if .CompressedFile⦘⦘
// _vfsgen_compressedFileInfo is a static definition of a gzip compressed file.
type _vfsgen_compressedFileInfo struct {
	name              string
	modTime           time.Time
	compressedContent []byte
	uncompressedSize  int64
}

func (f *_vfsgen_compressedFileInfo) Readdir(count int) ([]os.FileInfo, error) {
	return nil, fmt.Errorf("cannot Readdir from file %s", f.name)
}
func (f *_vfsgen_compressedFileInfo) Stat() (os.FileInfo, error) { return f, nil }

func (f *_vfsgen_compressedFileInfo) GzipBytes() []byte {
	return f.compressedContent
}

func (f *_vfsgen_compressedFileInfo) Name() string       { return f.name }
func (f *_vfsgen_compressedFileInfo) Size() int64        { return f.uncompressedSize }
func (f *_vfsgen_compressedFileInfo) Mode() os.FileMode  { return 0444 }
func (f *_vfsgen_compressedFileInfo) ModTime() time.Time { return f.modTime }
func (f *_vfsgen_compressedFileInfo) IsDir() bool        { return false }
func (f *_vfsgen_compressedFileInfo) Sys() interface{}   { return nil }

// _vfsgen_compressedFile is an opened compressedFile instance.
type _vfsgen_compressedFile struct {
	*_vfsgen_compressedFileInfo
	gr      *gzip.Reader
	grPos   int64 // Actual gr uncompressed position.
	seekPos int64 // Seek uncompressed position.
}

func (f *_vfsgen_compressedFile) Read(p []byte) (n int, err error) {
	if f.grPos > f.seekPos {
		// Rewind to beginning.
		err = f.gr.Reset(bytes.NewReader(f._vfsgen_compressedFileInfo.compressedContent))
		if err != nil {
			return 0, err
		}
		f.grPos = 0
	}
	if f.grPos < f.seekPos {
		// Fast-forward.
		_, err = io.ReadFull(f.gr, make([]byte, f.seekPos-f.grPos))
		if err != nil {
			return 0, err
		}
		f.grPos = f.seekPos
	}
	n, err = f.gr.Read(p)
	f.grPos += int64(n)
	f.seekPos = f.grPos
	return n, err
}
func (f *_vfsgen_compressedFile) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case os.SEEK_SET:
		f.seekPos = 0 + offset
	case os.SEEK_CUR:
		f.seekPos += offset
	case os.SEEK_END:
		f.seekPos = f._vfsgen_compressedFileInfo.uncompressedSize + offset
	}
	return f.seekPos, nil
}
func (f *_vfsgen_compressedFile) Close() error {
	return f.gr.Close()
}
⦗⦗else⦘⦘
// We already imported "compress/gzip", but ended up not using it. Avoid unused import error.
var _ = gzip.Reader
⦗⦗end⦘⦘⦗⦗if .File⦘⦘
// _vfsgen_fileInfo is a static definition of an uncompressed file (because it's not worth gzip compressing).
type _vfsgen_fileInfo struct {
	name    string
	modTime time.Time
	content []byte
}

func (f *_vfsgen_fileInfo) Readdir(count int) ([]os.FileInfo, error) {
	return nil, fmt.Errorf("cannot Readdir from file %s", f.name)
}
func (f *_vfsgen_fileInfo) Stat() (os.FileInfo, error) { return f, nil }

func (f *_vfsgen_fileInfo) NotWorthGzipCompressing() {} // TODO: Figure out a good interface to encode this information.

func (f *_vfsgen_fileInfo) Name() string       { return f.name }
func (f *_vfsgen_fileInfo) Size() int64        { return int64(len(f.content)) }
func (f *_vfsgen_fileInfo) Mode() os.FileMode  { return 0444 }
func (f *_vfsgen_fileInfo) ModTime() time.Time { return f.modTime }
func (f *_vfsgen_fileInfo) IsDir() bool        { return false }
func (f *_vfsgen_fileInfo) Sys() interface{}   { return nil }

// _vfsgen_file is an opened file instance.
type _vfsgen_file struct {
	*_vfsgen_fileInfo
	*bytes.Reader
}

func (f *_vfsgen_file) Close() error {
	return nil
}
⦗⦗end⦘⦘
// _vfsgen_dirInfo is a static definition of a directory.
type _vfsgen_dirInfo struct {
	name    string
	modTime time.Time
	entries []os.FileInfo
}

func (d *_vfsgen_dirInfo) Read([]byte) (int, error) {
	return 0, fmt.Errorf("cannot Read from directory %s", d.name)
}
func (d *_vfsgen_dirInfo) Close() error               { return nil }
func (d *_vfsgen_dirInfo) Stat() (os.FileInfo, error) { return d, nil }

func (d *_vfsgen_dirInfo) Name() string       { return d.name }
func (d *_vfsgen_dirInfo) Size() int64        { return 0 }
func (d *_vfsgen_dirInfo) Mode() os.FileMode  { return 0755 | os.ModeDir }
func (d *_vfsgen_dirInfo) ModTime() time.Time { return d.modTime }
func (d *_vfsgen_dirInfo) IsDir() bool        { return true }
func (d *_vfsgen_dirInfo) Sys() interface{}   { return nil }

// _vfsgen_dir is an opened dir instance.
type _vfsgen_dir struct {
	*_vfsgen_dirInfo
	pos int // Position within entries for Seek and Readdir.
}

func (d *_vfsgen_dir) Seek(offset int64, whence int) (int64, error) {
	if offset == 0 && whence == os.SEEK_SET {
		d.pos = 0
		return 0, nil
	}
	return 0, fmt.Errorf("unsupported Seek in directory %s", d._vfsgen_dirInfo.name)
}

func (d *_vfsgen_dir) Readdir(count int) ([]os.FileInfo, error) {
	if d.pos >= len(d._vfsgen_dirInfo.entries) && count > 0 {
		return nil, io.EOF
	}
	if count <= 0 || count > len(d._vfsgen_dirInfo.entries)-d.pos {
		count = len(d._vfsgen_dirInfo.entries) - d.pos
	}
	e := d._vfsgen_dirInfo.entries[d.pos : d.pos+count]
	d.pos += count
	return e, nil
}
`))
	return trailer.Execute(w, toc)
}
