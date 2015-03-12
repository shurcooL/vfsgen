// +build ignore

package main

import (
	"log"

	"github.com/shurcooL/go-bindata"
	"golang.org/x/tools/godoc/vfs/mapfs"
)

var inputFs = mapfs.New(map[string]string{
	"sample-file.txt":                "This file compresses well. Blaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaah.",
	"not-worth-compressing-file.txt": "Its normal contents are here.",
})

func main() {
	cfg := bindata.NewConfig()
	cfg.Input = inputFs
	cfg.Package = "bindata_test"
	cfg.Output = "./main_bindata_test.go"

	err := bindata.Translate(cfg)
	if err != nil {
		log.Fatalln(err)
	}
}
