package vfsgen

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	pathpkg "path"
	"sort"

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
	var toc []pathAsset
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

// findFiles recursively finds all the file paths in the given directory tree.
// They are added to the given map as keys. Values will be safe function names
// for each file, which will be used when generating the output code.
func findFiles(fs http.FileSystem, toc *[]pathAsset) error {
	walkFn := func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			log.Printf("can't stat file %s: %v\n", path, err)
			return nil
		}

		switch {
		case fi.IsDir():
			entries, err := readDirPaths(fs, path)
			if err != nil {
				return err
			}

			*toc = append(*toc, pathAsset{
				path: path,
				asset: &dir{
					name:    pathpkg.Base(path),
					entries: entries,
				},
			})

		case !fi.IsDir():
			*toc = append(*toc, pathAsset{
				path: path,
				asset: &compressedFile{
					name:             pathpkg.Base(path),
					uncompressedSize: fi.Size(),
				},
			})
		}

		return nil
	}

	err := vfsutil.Walk(fs, "/", walkFn)
	if err != nil {
		return err
	}

	return nil
}

type pathAsset struct {
	path  string
	asset interface{}
}

// writeAssets writes the code file.
func writeAssets(w io.Writer, c *Config, toc []pathAsset) error {
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
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, `type assetsFS map[string]interface{}

var AssetsFS http.FileSystem = func() assetsFS {
	assetsFS := assetsFS{
`)
	if err != nil {
		return err
	}

	for _, pathAsset := range toc {
		switch asset := pathAsset.asset.(type) {
		case *dir:
			_, err = fmt.Fprintf(w, "\t\t%q: &dir{\n", pathAsset.path)
			if err != nil {
				return err
			}
			fmt.Fprintf(w, "\t\t\tname:    %q,\n", asset.name)
			fmt.Fprintf(w, "\t\t\tmodTime: time.Time{},\n")
			fmt.Fprintf(w, "\t\t},\n")
		case *compressedFile:
			_, err = fmt.Fprintf(w, "\t\t%q: &compressedFile{\n", pathAsset.path)
			if err != nil {
				return err
			}
			fmt.Fprintf(w, "\t\t\tname:              %q,\n", asset.name)
			fmt.Fprintf(w, "\t\t\tcompressedContent: []byte(\"")
			f, _ := c.Input.Open(pathAsset.path)
			sw := &StringWriter{Writer: w}
			gz := gzip.NewWriter(sw)
			io.Copy(gz, f)
			gz.Close()
			f.Close()
			fmt.Fprintf(w, "\"),\n")
			fmt.Fprintf(w, "\t\t\tuncompressedSize:  %d,\n", asset.uncompressedSize)
			fmt.Fprintf(w, "\t\t\tmodTime:           time.Time{},\n")
			fmt.Fprintf(w, "\t\t},\n")
		}
	}

	_, err = fmt.Fprintf(w, "\t}\n\n")
	if err != nil {
		return err
	}

	for _, pathAsset := range toc {
		switch asset := pathAsset.asset.(type) {
		case *dir:
			fmt.Fprintf(w, "\tassetsFS[%q].(*dir).entries = []os.FileInfo{\n", pathAsset.path)
			for _, entry := range asset.entries {
				fmt.Fprintf(w, "\t\tassetsFS[%q].(os.FileInfo),\n", entry)
			}
			fmt.Fprintf(w, "\t}\n")
		}
	}

	_, err = fmt.Fprintf(w, "\n\treturn assetsFS\n}()\n")
	if err != nil {
		return err
	}

	return nil
}

/*// writeAsset write a entry for the given asset.
// An entry is a function which embeds and returns
// the file's byte content.
func writeAsset(w io.Writer, c *Config, asset interface{}) error {
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

	_, err = fmt.Fprintf(w, "\t\t%q: ...,\n", asset.Name)
	if err != nil {
		return err
	}

	return nil
}*/

/*func writeCompressedAsset(w io.Writer, asset *Asset, r io.Reader) (int64, error) {
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
}*/

func writeVFS(w io.Writer) error {
	_, err := fmt.Fprint(w, `
func (fs assetsFS) Open(path string) (http.File, error) {
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
