package vfsgen

import (
	"go/build"
	"log"
	"net/http"
)

// DevDir returns a http.FileSystem reading from the first directory found
// inside $GOPATH/src with the given name. That is:
//
//  DevDir("github.com/my/project/assets")
//
// Would return an http.FileSystem reading from
//
//  $GOPATH/src/github.com/my/project/assets
//
// Where $GOPATH is the first $GOPATH entry that was found to contain that
// directory. If the directory does not exist, the program exits fatally.
func DevDir(path string) http.FileSystem {
	pkg, err := build.Import(path, "", build.FindOnly)
	if err != nil {
		log.Fatal(err)
	}
	return http.FileSystem(http.Dir(pkg.Dir))
}
