package bindata_test

import (
	"fmt"
	"log"
	"os"

	"github.com/shurcooL/go/vfs_util"
	"github.com/shurcooL/go/vfsfs"
	"golang.org/x/tools/godoc/vfs"
)

//go:generate go run main_test_generate.go

func Example() {
	var fs vfs.FileSystem = vfsfs.New(&AssetFS{Asset: Asset, AssetDir: AssetDir})

	walkFn := func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			log.Printf("can't stat file %s: %v\n", path, err)
			return nil
		}
		fmt.Println(path)
		if !fi.IsDir() {
			b, err := vfs.ReadFile(fs, path)
			fmt.Printf("%q %v\n", string(b), err)
		}
		return nil
	}

	err := vfs_util.Walk(fs, "/", walkFn)
	if err != nil {
		panic(err)
	}

	// Output:
	// /
	// /sample-file.txt
	// "Its normal contents are here." <nil>
}
