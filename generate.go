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
	err := header_compressed_nomemcopy(w)
	if err != nil {
		return err
	}
	return header_common(w)
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

	compressedSize, err := compressed_nomemcopy(w, asset, fd)
	if err != nil {
		return err
	}
	return asset_common(w, c, asset, compressedSize)
}

func header_compressed_nomemcopy(w io.Writer) error {
	_, err := fmt.Fprintf(w, `import (
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
		return nil, fmt.Errorf("Read %%q: %%v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %%q: %%v", name, err)
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

`)
	return err
}

func header_common(w io.Writer) error {
	_, err := fmt.Fprintf(w, `type asset struct {
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

func compressed_nomemcopy(w io.Writer, asset *Asset, r io.Reader) (int64, error) {
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

func asset_common(w io.Writer, c *Config, asset *Asset, compressedSize int64) error {
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
