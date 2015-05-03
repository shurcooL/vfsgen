package vfsgen_test

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/shurcooL/go/vfs/httpfs/vfsutil"
)

//go:generate go run main_test_generate.go

func Example() {
	var fs http.FileSystem = AssetsFS

	walkFn := func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			log.Printf("can't stat file %s: %v\n", path, err)
			return nil
		}

		fmt.Println(path)
		if fi.IsDir() {
			return nil
		}

		b, err := vfsutil.ReadFile(fs, path)
		fmt.Printf("%q %v\n", string(b), err)
		return nil
	}

	err := vfsutil.Walk(fs, "/", walkFn)
	if err != nil {
		panic(err)
	}

	// Output:
	// /
	// /folderA
	// /folderA/file1.txt
	// "Stuff." <nil>
	// /folderA/file2.txt
	// "Stuff." <nil>
	// /folderB
	// /folderB/folderC
	// /folderB/folderC/file3.txt
	// "Stuff." <nil>
	// /not-worth-compressing-file.txt
	// "Its normal contents are here." <nil>
	// /sample-file.txt
	// "This file compresses well. Blaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaah!" <nil>
}

type gzipByter interface {
	GzipBytes() []byte
}

func ExampleCompressed() {
	// Compressed file system.
	var fs http.FileSystem = AssetsFS

	walkFn := func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			log.Printf("can't stat file %s: %v\n", path, err)
			return nil
		}

		fmt.Println(path)
		if fi.IsDir() {
			return nil
		}

		f, err := fs.Open(path)
		if err != nil {
			fmt.Printf("fs.Open(%q): %v\n", path, err)
			return nil
		}
		defer f.Close()

		b, err := ioutil.ReadAll(f)
		fmt.Printf("%q %v\n", string(b), err)

		if gzipFile, ok := f.(gzipByter); ok {
			b := gzipFile.GzipBytes()
			fmt.Printf("%q\n", string(b))
		} else {
			fmt.Println("<not compressed>")
		}
		return nil
	}

	err := vfsutil.Walk(fs, "/", walkFn)
	if err != nil {
		panic(err)
	}

	// Output:
	// /
	// /folderA
	// /folderA/file1.txt
	// "Stuff." <nil>
	// "\x1f\x8b\b\x00\x00\tn\x88\x00\xff\n.)MK\xd3\x03\x04\x00\x00\xff\xff'\xbb@\xc8\x06\x00\x00\x00"
	// /folderA/file2.txt
	// "Stuff." <nil>
	// "\x1f\x8b\b\x00\x00\tn\x88\x00\xff\n.)MK\xd3\x03\x04\x00\x00\xff\xff'\xbb@\xc8\x06\x00\x00\x00"
	// /folderB
	// /folderB/folderC
	// /folderB/folderC/file3.txt
	// "Stuff." <nil>
	// "\x1f\x8b\b\x00\x00\tn\x88\x00\xff\n.)MK\xd3\x03\x04\x00\x00\xff\xff'\xbb@\xc8\x06\x00\x00\x00"
	// /not-worth-compressing-file.txt
	// "Its normal contents are here." <nil>
	// "\x1f\x8b\b\x00\x00\tn\x88\x00\xff\xf2,)V\xc8\xcb/\xcaM\xccQH\xce\xcf+I\xcd\x03\xf2\x13\x8bR\x152R\x8bR\xf5\x00\x01\x00\x00\xff\xff\xdc\xc7\xff\x13\x1d\x00\x00\x00"
	// /sample-file.txt
	// "This file compresses well. Blaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaah!" <nil>
	// "\x1f\x8b\b\x00\x00\tn\x88\x00\xff\n\xc9\xc8,VH\xcb\xccIUH\xce\xcf-(J-.N-V(O\xcd\xc9\xd1Sp\xcaI\x1c\xd4 C\x11\x10\x00\x00\xff\xff\xe7G\x81:\xbd\x00\x00\x00"
}

func ExampleReadTwoOpenedFiles() {
	var fs http.FileSystem = AssetsFS

	f0, err := fs.Open("/sample-file.txt")
	if err != nil {
		panic(err)
	}
	defer f0.Close()
	f1, err := fs.Open("/sample-file.txt")
	if err != nil {
		panic(err)
	}
	defer f1.Close()

	_, err = io.CopyN(os.Stdout, f0, 9)
	if err != nil {
		panic(err)
	}
	_, err = io.CopyN(os.Stdout, f1, 9)
	if err != nil {
		panic(err)
	}

	// Output:
	// This fileThis file
}

func ExampleModTime() {
	var fs http.FileSystem = AssetsFS

	f, err := fs.Open("/sample-file.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		panic(err)
	}

	fmt.Println(fi.ModTime())

	// Output:
	// 0001-01-01 00:00:00 +0000 UTC
}

func ExampleSeek() {
	var fs http.FileSystem = AssetsFS

	f, err := fs.Open("/sample-file.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = io.CopyN(os.Stdout, f, 5)
	if err != nil {
		panic(err)
	}
	_, err = f.Seek(22, os.SEEK_CUR)
	if err != nil {
		panic(err)
	}
	_, err = io.CopyN(os.Stdout, f, 10)
	if err != nil {
		panic(err)
	}
	fmt.Print("...")
	_, err = f.Seek(-4, os.SEEK_END)
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(os.Stdout, f)
	if err != nil {
		panic(err)
	}
	_, err = f.Seek(3, os.SEEK_SET)
	if err != nil {
		panic(err)
	}
	_, err = f.Seek(1, os.SEEK_CUR)
	if err != nil {
		panic(err)
	}
	_, err = io.CopyN(os.Stdout, f, 22)
	if err != nil {
		panic(err)
	}

	// Output:
	// This Blaaaaaaaa...aah! file compresses well.
}
