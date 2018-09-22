//go:generate go run github.com/shurcooL/vfsgen/cmd/vfsgendev -source="example.com/foo/example".Assets
package main

import (
	"fmt"
	"net/http"

	"github.com/shurcooL/httpgzip"
)

func main() {
	var Assets http.FileSystem = http.Dir("./assets")

	file1, err := Assets.Open("/hello.txt")
	if nil != err {
		fmt.Println("Problemo:", err)
	}
	_, ok := file1.(httpgzip.GzipByter)
	if !ok {
		fmt.Println("hello.txt was not gzipped")
	}
	defer file1.Close()

	file2, _ := Assets.Open("/alphabet-alphabet.txt")
	_, ok = file2.(httpgzip.GzipByter)
	if ok {
		fmt.Println("alphabet-alphabet was gzipped")
	}
	defer file2.Close()

}
