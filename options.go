package vfsgen

import "fmt"

// Options specifies options for vfsgen code generation.
type Options struct {
	// Filename is the output Go file filename (including extension) for the generated code.
	// If left empty, it defaults to "{{.VariableName}}_vfsdata.go".
	Filename string

	// PackageName is the name of the package in the generated code.
	// If left empty, it defaults to "main".
	PackageName string

	// BuildTags are the optional build tags in the generated code.
	// The build tags syntax is specified by the go tool.
	BuildTags string

	// VariableName is the name of the http.FileSystem variable in the generated code.
	// If left empty, it defaults to "assets".
	VariableName string
}

// fillMissing sets default values for mandatory options that are left empty.
func (opt *Options) fillMissing() {
	if opt.PackageName == "" {
		opt.PackageName = "main"
	}
	if opt.VariableName == "" {
		opt.VariableName = "assets"
	}
	if opt.Filename == "" {
		opt.Filename = fmt.Sprintf("%s_vfsdata.go", opt.VariableName)
	}
}
