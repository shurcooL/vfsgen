// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0/

package bindata

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"golang.org/x/tools/godoc/vfs"
)

// Config defines a set of options for the asset conversion.
type Config struct {
	// Name of the package to use. Defaults to 'main'.
	Package string

	// Tags specify a set of optional build tags, which should be
	// included in the generated output. The tags are appended to a
	// `// +build` line in the beginning of the output file
	// and must follow the build tags syntax specified by the go tool.
	Tags string

	// Input ...
	Input vfs.FileSystem

	// Output defines the output file for the generated code.
	// If left empty, this defaults to 'bindata.go' in the current
	// working directory.
	Output string

	// Prefix defines a path prefix which should be stripped from all
	// file names when generating the keys in the table of contents.
	// For example, running without the `-prefix` flag, we get:
	//
	// 	$ go-bindata /path/to/templates
	// 	go_bindata["/path/to/templates/foo.html"] = _path_to_templates_foo_html
	//
	// Running with the `-prefix` flag, we get:
	//
	// 	$ go-bindata -prefix "/path/to/" /path/to/templates/foo.html
	// 	go_bindata["templates/foo.html"] = templates_foo_html
	Prefix string

	// Ignores any filenames matching the regex pattern specified, e.g.
	// path/to/file.ext will ignore only that file, or \\.gitignore
	// will match any .gitignore file.
	//
	// This parameter can be provided multiple times.
	Ignore []*regexp.Regexp
}

// NewConfig returns a default configuration struct.
func NewConfig() *Config {
	return &Config{
		Package: "main",
		Output:  "./bindata.go",
		Ignore:  make([]*regexp.Regexp, 0),
	}
}

// validate ensures the config has sane values.
// Part of which means checking if certain file/directory paths exist.
func (c *Config) validate() error {
	if len(c.Package) == 0 {
		return fmt.Errorf("Missing package name")
	}

	_, err := c.Input.Lstat("/")
	if err != nil {
		return fmt.Errorf("Failed to stat input root: %v", err)
	}

	if len(c.Output) == 0 {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("Unable to determine current working directory.")
		}

		c.Output = filepath.Join(cwd, "bindata.go")
	}

	stat, err := os.Lstat(c.Output)
	if err != nil {
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
	}

	if stat != nil && stat.IsDir() {
		return fmt.Errorf("Output path is a directory.")
	}

	return nil
}
