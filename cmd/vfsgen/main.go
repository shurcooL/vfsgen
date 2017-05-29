// vfsgen generates Go code that statically implements an http.FileSystem for a
// given directory.
//
// Installation
//
//  go get -u github.com/shurcooL/vfsgen/cmd/vfsgen
//
// Basic Usage
//
// To generate a assets_vfsdata.go file containing everything in mydir/ for your
// package, you can run:
//
//  vfsgen -pkg=mypkg -dir=mydir/
//
// To see other configuration options run 'vfsgen -h'.
//
// Filters
//
// vfsgen supports various combinations of filter expressions. For example to
// exclude all .go source files:
//
//  vfsgen -filter='Extensions(".go", ".html")'
//
// Also same as the above is to use the Combine function to combine multiple
// filters together (effectively an OR operator):
//
//  vfsgen -filter='Combine(Extensions(".go"), Extensions(".html"))'
//
// Or exclude anything that is not a .html file using the Not function:
//
//  vfsgen -filter='Not(Extensions(".html"))'
//
// More complex situations can be derived from these examples, obviously. Also
// note that the single quotes '' are only needed for execution in Bash,
// go:generate directives do not need these.
//
// Go Generate
//
// One common usage pattern is to use vfsgen as part of a go:generate directive,
// for example add a generate.go to your package:
//
//  //go:generate vfsgen -pkg=assets -tags=!dev -filter=Extensions(".go")
//
//  package assets
//
// And thus running:
//
//  go generate github.com/my/pkg/assets
//
// Would produce a assets_vfsdata.go with any non-go file in the assets
// directory.
//
// Development
//
// Because we told vfsgen to not build this file when the dev build tag is
// present, you could add another file for development to your project at
// github.com/my/pkg/assets/assets_dev.go like:
//
//  // +build dev
//
//  package assets
//
//  import "github.com/shurcooL/vfsgen"
//
//  var Assets = vfsgen.DevDir("github.com/my/pkg/assets")
//
// Now just build your project with the development build tag:
//
//  go install -tags=dev github.com/my/pkg
//
// And files will be read directly from the filesystem instead of from the
// packaged assets.
//
// Use as a package
//
// One different mode of execution is using vfsgen as a Go package. This means
// you specify the options in Go syntax rather than as CLI parameters. We
// suggest using this mode if you find yourself with very long / hard to read
// CLI flag combinations. For more details see the example at:
//
// https://godoc.org/github.com/shurcooL/vfsgen
//
package main

import (
	"flag"
	"log"
	"net/http"
	"strings"

	"github.com/shurcooL/httpfs/filter"
	"github.com/shurcooL/vfsgen"
)

var (
	filename   = flag.String("filename", "", "output filename (defaults to \"{{toLower .VariableName}}_vfsdata.go\")")
	pkgName    = flag.String("pkg", "main", "package name to emit")
	buildTags  = flag.String("tags", "", "build tags to emit")
	varName    = flag.String("var", "Assets", "variable name of http.Dir to emit")
	dir        = flag.String("dir", ".", "directory to generate VFS from")
	filterFlag = flag.String("filter", "", "filters to apply")
)

func main() {
	log.SetFlags(0)
	flag.Parse()

	fs := http.FileSystem(http.Dir(*dir))

	// Install the filter, if needed.
	if *filterFlag != "" {
		ignore, err, node := filter.Parse(*filterFlag)
		if err != nil {
			log.Printf("filter: inside expression %s\n", *filterFlag)
			if node != nil {
				start := strings.Repeat(" ", int(node.Pos())-2)
				log.Fatalf("filter:                   %v ^ %v\n", start, err)
			} else {
				log.Fatalln("filter:", err)
			}
		}
		fs = filter.New(fs, ignore)
	}
	err := vfsgen.Generate(fs, vfsgen.Options{
		Filename:     *filename,
		PackageName:  *pkgName,
		BuildTags:    *buildTags,
		VariableName: *varName,
	})
	if err != nil {
		log.Fatalln(err)
	}
}
