package vfsgen

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/shurcooL/go/vfs/httpfs/vfsutil"
)

// Config defines a set of options for the asset conversion.
type Config struct {
	// Input is the filesystem that contains input assets to be converted.
	Input http.FileSystem

	// THINKING:
	// Input
	// OutputFile
	// OutputPackageTags
	// OutputPackageName
	// OutputVariableName

	// Output defines the output file for the generated code.
	// If left empty, this defaults to "./{{.OutputName}}_vfsdata.go".
	Output string

	// Tags specify a set of optional build tags, which should be
	// included in the generated output. The tags are appended to a
	// `// +build` line in the beginning of the output file
	// and must follow the build tags syntax specified by the go tool.
	Tags string

	// Name of the package to use. Defaults to 'main'.
	Package string

	// OutputName defines the output filesystem variable name.
	// If left empty, this defaults to "assets".
	OutputName string
}

// validate ensures the config has sane values.
// Part of which means checking if certain file/directory paths exist.
func (c *Config) validate() error {
	if c.Package == "" {
		c.Package = "main"
	}

	_, err := vfsutil.Stat(c.Input, "/")
	if err != nil {
		return fmt.Errorf("Failed to stat input root: %v", err)
	}

	if c.OutputName == "" {
		c.OutputName = "assets"
	}

	if c.Output == "" {
		c.Output = fmt.Sprintf("./%s_vfsdata.go", c.OutputName)
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
