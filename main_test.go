package bindata_test

import (
	"fmt"
	"log"
	"os"

	"github.com/shurcooL/go/vfs_util"
	"golang.org/x/tools/godoc/vfs"
)

//go:generate go run main_test_generate.go

func Example() {
	var fs vfs.FileSystem = AssetsFs

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
	// /not-worth-compressing-file.txt
	// "Its normal contents are here." <nil>
	// /sample-file.txt
	// "This file compresses well. Blaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaah." <nil>
}

func ExampleCompressed() {
	// Compressed file system.
	var fs vfs.FileSystem = AssetsFs

	walkFn := func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			log.Printf("can't stat file %s: %v\n", path, err)
			return nil
		}
		fmt.Println(path)
		if !fi.IsDir() {
			// if func (f *AssetFile) CompressedBytes() ([]byte, error) {...
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
	// /not-worth-compressing-file.txt
	// "\x1f\x8b\b\x00\x00\tn\x88\x00\xff\xf2,)V\xc8\xcb/\xcaM\xccQH\xce\xcf+I\xcd\x03\xf2\x13\x8bR\x152R\x8bR\xf5\x00\x01\x00\x00\xff\xff\xdc\xc7\xff\x13\x1d\x00\x00\x00" <nil>
	// /sample-file.txt
	// "\x1f\x8b\b\x00\x00\tn\x88\x00\xff\n\xc9\xc8,VH\xcb\xccIUH\xce\xcf-(J-.N-V(O\xcd\xc9\xd1Sp\xcaI\x1c\xd4 C\x0f\x10\x00\x00\xff\xffvZ>\xaa\xbd\x00\x00\x00" <nil>
}
