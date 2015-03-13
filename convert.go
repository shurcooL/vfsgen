// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0/

package bindata

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/shurcooL/go/vfs/godocfs/vfsutil"
	"golang.org/x/tools/godoc/vfs"
)

// Translate reads assets from an input directory, converts them
// to Go code and writes new files to the output specified
// in the given configuration.
func Translate(c *Config) error {
	var toc []Asset

	// Ensure our configuration has sane values.
	err := c.validate()
	if err != nil {
		return err
	}

	var knownFuncs = make(map[string]int)
	// Locate all the assets.
	err = findFiles(c.Input, &toc, c.Ignore, knownFuncs)
	if err != nil {
		return err
	}

	// Create output file.
	fd, err := os.Create(c.Output)
	if err != nil {
		return err
	}
	defer fd.Close()

	// Create a buffered writer for better performance.
	bfd := bufio.NewWriter(fd)
	defer bfd.Flush()

	// Write generated disclaimer.
	_, err = fmt.Fprintf(bfd, "// generated via `go generate`; do not edit\n\n")
	if err != nil {
		return err
	}

	// Write build tags, if applicable.
	if c.Tags != "" {
		_, err = fmt.Fprintf(bfd, "// +build %s\n\n", c.Tags)
		if err != nil {
			return err
		}
	}

	// Write package declaration.
	_, err = fmt.Fprintf(bfd, "package %s\n\n", c.Package)
	if err != nil {
		return err
	}

	// Write assets.
	if err := writeAssets(bfd, c, toc); err != nil {
		return err
	}

	// Write table of contents.
	if err := writeTOC(bfd, toc); err != nil {
		return err
	}

	// Write hierarchical tree of assets.
	if err := writeTOCTree(bfd, toc); err != nil {
		return err
	}

	// Write virtual file system.
	if err := writeVFS(bfd); err != nil {
		return err
	}

	return nil
}

// Implement sort.Interface for []os.FileInfo based on Name()
type ByName []os.FileInfo

func (v ByName) Len() int           { return len(v) }
func (v ByName) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v ByName) Less(i, j int) bool { return v[i].Name() < v[j].Name() }

// findFiles recursively finds all the file paths in the given directory tree.
// They are added to the given map as keys. Values will be safe function names
// for each file, which will be used when generating the output code.
func findFiles(dir vfs.FileSystem, toc *[]Asset, ignore []*regexp.Regexp, knownFuncs map[string]int) error {
	walkFn := func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			log.Printf("can't stat file %s: %v\n", path, err)
			return nil
		}

		var asset Asset
		asset.Path = path
		asset.Name = path

		ignoring := false
		for _, re := range ignore {
			if re.MatchString(asset.Path) {
				ignoring = true
				break
			}
		}
		if ignoring {
			return nil
		}

		/*if fi.IsDir() {
			if recursive {
				return nil
			} else {
				return filepath.SkipDir
			}
		}*/
		if fi.IsDir() {
			return nil
		}

		// If we have a leading slash, get rid of it.
		if len(asset.Name) > 0 && asset.Name[0] == '/' {
			asset.Name = asset.Name[1:]
		}

		// This shouldn't happen.
		if len(asset.Name) == 0 {
			return fmt.Errorf("Invalid file: %v", asset.Path)
		}

		asset.Func = safeFunctionName(asset.Name, knownFuncs)
		asset.Path, _ = filepath.Abs(asset.Path)
		*toc = append(*toc, asset)

		return nil
	}

	err := vfsutil.Walk(dir, "/", walkFn)
	if err != nil {
		return err
	}

	return nil
}

var regFuncName = regexp.MustCompile(`[^a-zA-Z0-9_]`)
