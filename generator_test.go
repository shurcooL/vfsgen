package vfsgen_test

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/shurcooL/httpfs/union"
	"github.com/shurcooL/vfsgen"
	"golang.org/x/tools/godoc/vfs/httpfs"
	"golang.org/x/tools/godoc/vfs/mapfs"
)

// This code will generate an assets_vfsdata.go file with
// `var assets http.FileSystem = ...`
// that statically implements the contents of "assets" directory.
//
// vfsgen is great to use with go generate directives. This code can go in an assets_gen.go file, which can
// then be invoked via "//go:generate go run assets_gen.go". The input virtual filesystem can read directly
// from disk, or it can be more involved.
func Example() {
	var fs http.FileSystem = http.Dir("assets")

	err := vfsgen.Generate(fs, vfsgen.Options{})
	if err != nil {
		log.Fatalln(err)
	}
}

// Verify that all possible combinations of {non-compressed,compressed} files build
// successfully, and have no gofmt issues.
func TestGenerate_buildAndGofmt(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "vfsgen_test_")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Fatal(err)
		}
	}()

	tests := []struct {
		filename  string
		fs        http.FileSystem
		wantError func(error) bool // Nil function means want nil error.
	}{
		{
			// Empty.
			filename: "empty.go",
			fs:       union.New(nil),
		},
		{
			// Test that vfsgen.Generate returns an error when there is
			// an error reading from the input filesystem.
			filename:  "notexist.go",
			fs:        http.Dir("notexist"),
			wantError: os.IsNotExist,
		},
		{
			// No compressed files.
			filename: "nocompressed.go",
			fs: httpfs.New(mapfs.New(map[string]string{
				"not-compressable-file.txt": "Not compressable.",
			})),
		},
		{
			// Only compressed files.
			filename: "onlycompressed.go",
			fs: httpfs.New(mapfs.New(map[string]string{
				"compressable-file.txt": "This text compresses easily. " + strings.Repeat(" Go!", 128),
			})),
		},
		{
			// Both non-compressed and compressed files.
			filename: "both.go",
			fs: httpfs.New(mapfs.New(map[string]string{
				"not-compressable-file.txt": "Not compressable.",
				"compressable-file.txt":     "This text compresses easily. " + strings.Repeat(" Go!", 128),
			})),
		},
	}

	for _, test := range tests {
		filename := filepath.Join(tempDir, test.filename)

		err := vfsgen.Generate(test.fs, vfsgen.Options{
			Filename:    filename,
			PackageName: "test",
		})
		switch {
		case test.wantError == nil && err != nil:
			t.Fatalf("%s: vfsgen.Generate returned non-nil error: %v", test.filename, err)
		case test.wantError != nil && !test.wantError(err):
			t.Fatalf("%s: vfsgen.Generate returned wrong error: %v", test.filename, err)
		}
		if test.wantError != nil {
			continue
		}

		if out, err := exec.Command("go", "build", filename).CombinedOutput(); err != nil {
			t.Errorf("err: %v\nout: %s", err, out)
		}
		if out, err := exec.Command("gofmt", "-d", "-s", filename).Output(); err != nil || len(out) != 0 {
			t.Errorf("gofmt issue\nerr: %v\nout: %s", err, out)
		}
	}
}
