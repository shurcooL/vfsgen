package vfsgen

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/shurcooL/go/vfs/httpfs/vfsutil"
)

// NewConfig returns a default configuration struct.
func NewConfig() *Config {
	return &Config{
		Package: "main",
		Output:  "./vfsdata.go",
		//OutputName: "AssetsFs",
	}
}

// Config defines a set of options for the asset conversion.
type Config struct {
	// Input is the filesystem that contains input assets to be converted.
	Input http.FileSystem

	// Name of the package to use. Defaults to 'main'.
	Package string

	// Tags specify a set of optional build tags, which should be
	// included in the generated output. The tags are appended to a
	// `// +build` line in the beginning of the output file
	// and must follow the build tags syntax specified by the go tool.
	Tags string

	// Output defines the output file for the generated code.
	// If left empty, this defaults to "vfsdata.go" in the current
	// working directory.
	Output string

	// OutputName defines the output filesystem variable name.
	// If left empty, this defaults to "AssetsFs".
	//OutputName string
}

// validate ensures the config has sane values.
// Part of which means checking if certain file/directory paths exist.
func (c *Config) validate() error {
	if len(c.Package) == 0 {
		return fmt.Errorf("Missing package name")
	}

	_, err := vfsutil.Stat(c.Input, "/")
	if err != nil {
		return fmt.Errorf("Failed to stat input root: %v", err)
	}

	if len(c.Output) == 0 {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("Unable to determine current working directory.")
		}

		c.Output = filepath.Join(cwd, "vfsdata.go")
	}

	switch stat, err := os.Lstat(c.Output); {
	case err != nil:
		if !os.IsNotExist(err) {
			return fmt.Errorf("Output path: %v", err)
		}

		// File does not exist. This is fine, just make
		// sure the directory it is to be in exists.
		dir, _ := filepath.Split(c.Output)
		if dir != "" {
			err = os.MkdirAll(dir, 0744)

			if err != nil {
				return fmt.Errorf("Create output directory: %v", err)
			}
		}
	case stat.IsDir():
		return fmt.Errorf("Output path is a directory.")
	}

	return nil
}
