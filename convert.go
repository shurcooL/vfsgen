package vfsgen

import (
	"bufio"
	"fmt"
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
