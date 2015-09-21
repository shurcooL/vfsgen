package vfsgen

import "fmt"

// Options specifies options for vfsgen code generation.
type Options struct {
	// Filename is the output filename (including extension) for the generated code.
	// If left empty, this defaults to "{{.VariableName}}_vfsdata.go".
	Filename string

	// Package is the name of the package in the generated code.
	// If left empty, this defaults to "main".
	Package string

	// Tags is the optional build tags in the generated code.
	// Tags must follow the build tags syntax specified by the go tool.
	Tags string

	// VariableName is the name of the http.FileSystem variable in the generated code.
	// If left empty, this defaults to "assets".
	VariableName string
}

// fillMissing sets default values for mandatory options that are left blank.
func (opt *Options) fillMissing() {
	if opt.Package == "" {
		opt.Package = "main"
	}
	if opt.VariableName == "" {
		opt.VariableName = "assets"
	}
	if opt.Filename == "" {
		opt.Filename = fmt.Sprintf("%s_vfsdata.go", opt.VariableName)
	}
}
