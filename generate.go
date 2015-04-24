package vfsgen

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/shurcooL/go/vfs/httpfs/vfsutil"
)

// Translate reads assets from an input directory, converts them
// to Go code and writes new files to the output specified
// in the given configuration.
func Translate(c *Config) error {
	// Ensure our configuration has sane values.
	err := c.validate()
	if err != nil {
		return err
	}

	// Locate all the assets.
	var toc []Asset
	err = findFiles(c.Input, &toc)
	if err != nil {
		return err
	}

	// Create output file.
	f, err := os.Create(c.Output)
	if err != nil {
		return err
	}
	defer f.Close()

	// Create a buffered writer for better performance.
	buf := bufio.NewWriter(f)
	defer buf.Flush()

	// Write generated disclaimer.
	_, err = fmt.Fprintf(buf, "// generated via `go generate`; do not edit\n\n")
	if err != nil {
		return err
	}

	// Write build tags, if applicable.
	if c.Tags != "" {
		_, err = fmt.Fprintf(buf, "// +build %s\n\n", c.Tags)
		if err != nil {
			return err
		}
	}

	// Write package declaration.
	_, err = fmt.Fprintf(buf, "package %s\n\n", c.Package)
	if err != nil {
		return err
	}

	// Write assets.
	err = writeAssets(buf, c, toc)
	if err != nil {
		return err
	}

	// Write virtual file system.
	err = writeVFS(buf)
	if err != nil {
		return err
	}

	return nil
}

// TODO.
//
// Asset holds information about a single asset to be processed.
type Asset struct {
	Path string // Full file path.
	Name string // Key used in TOC -- name by which asset is referenced.
	Func string // Function name for the procedure returning the asset contents.
}

// findFiles recursively finds all the file paths in the given directory tree.
// They are added to the given map as keys. Values will be safe function names
// for each file, which will be used when generating the output code.
func findFiles(fs http.FileSystem, toc *[]Asset) error {
	walkFn := func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			log.Printf("can't stat file %s: %v\n", path, err)
			return nil
		}

		var asset Asset
		asset.Path = path
		asset.Name = path

		if fi.IsDir() {
			return nil
		}

		// If we have a leading slash, get rid of it.
		asset.Name = strings.TrimPrefix(asset.Name, "/")

		// This shouldn't happen.
		if len(asset.Name) == 0 {
			return fmt.Errorf("Invalid file: %v", asset.Path)
		}

		//asset.Func = safeFunctionName(asset.Name, knownFuncs)
		//*toc = append(*toc, asset)

		return nil
	}

	err := vfsutil.Walk(fs, "/", walkFn)
	if err != nil {
		return err
	}

	return nil
}

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
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

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
`)
	return err
}
