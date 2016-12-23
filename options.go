package vfsgen

import (
	"fmt"
	"strings"
)

// Options for vfsgen code generation.
type Options struct {
	// DataFilename is the filename of the generated VFS data Go code (including extension).
	// If left empty, it defaults to "{{toLower .VariableName}}_vfsdata.go".
	DataFilename string

	// CommonFilename is the filename of the generated common vfsgen Go code (including extension).
	// If left empty, it defaults to "vfsgencommon.go".
	CommonFilename string

	// PackageName is the name of the package in the generated code.
	// If left empty, it defaults to "main".
	PackageName string

	// BuildTags are the optional build tags in the generated code.
	// The build tags syntax is specified by the go tool.
	BuildTags string

	// VariableName is the name of the http.FileSystem variable in the generated code.
	// If left empty, it defaults to "assets".
	VariableName string

	// VariableComment is the comment of the http.FileSystem variable in the generated code.
	// If left empty, it defaults to "{{.VariableName}} statically implements the virtual filesystem provided to vfsgen.".
	VariableComment string
}

// fillMissing sets default values for mandatory options that are left empty.
func (opt *Options) fillMissing() {
	if opt.PackageName == "" {
		opt.PackageName = "main"
	}
	if opt.CommonFilename == "" {
		opt.CommonFilename = "vfsgencommon.go"
	}
	if opt.VariableName == "" {
		opt.VariableName = "assets"
	}
	if opt.DataFilename == "" {
		opt.DataFilename = fmt.Sprintf("%s_vfsdata.go", strings.ToLower(opt.VariableName))
	}
	if opt.VariableComment == "" {
		opt.VariableComment = fmt.Sprintf("%s statically implements the virtual filesystem provided to vfsgen.", opt.VariableName)
	}
}
