package vfsgen_test

import (
	"log"
	"net/http"

	"github.com/shurcooL/vfsgen"
)

// This code will generate an assets_vfsdata.go file with
// `var assets http.FileSystem = ...`
// that statically implements the contents of "assets" directory.
//
// Then, in your program, you can use `assets` as any other `http.FileSystem`, for example:
//
// 	file, err := assets.Open("/some/file.txt")
// 	if err != nil { ... }
// 	defer file.Close()
//
// 	http.Handle("/assets/", http.FileServer(assets))
//
// vfsgen is great to use with go generate directives. This code can go in an assets_gen.go file, which can
// then be invoked via "//go:generate go run assets_gen.go". The input virtual filesystem can read directly
// from disk, or it can be something more involved.
func Example() {
	var fs http.FileSystem = http.Dir("assets")

	err := vfsgen.Generate(fs, vfsgen.Options{})
	if err != nil {
		log.Fatalln(err)
	}
}
