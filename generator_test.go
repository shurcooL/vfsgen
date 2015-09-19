package vfsgen_test

import (
	"log"
	"net/http"

	"github.com/shurcooL/vfsgen"
)

// This code will generate a assets_vfsdata.go file that statically implements the contents of "assets" directory.
//
// It is typically meant to be executed via go generate directives. This code can go in an assets_gen.go file,
// which can then be invoked via "//go:generate go run assets_gen.go". The input virtual filesystem can read
// directly from disk, or it can be something more involved.
func Example() {
	var fs http.FileSystem = http.Dir("assets")

	config := vfsgen.Config{
		Input: fs,
	}

	err := vfsgen.Generate(config)
	if err != nil {
		log.Fatalln(err)
	}
}
