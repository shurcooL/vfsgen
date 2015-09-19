# vfsgen [![Build Status](https://travis-ci.org/shurcooL/vfsgen.svg?branch=master)](https://travis-ci.org/shurcooL/vfsgen) [![GoDoc](https://godoc.org/github.com/shurcooL/vfsgen?status.svg)](https://godoc.org/github.com/shurcooL/vfsgen)

Package vfsgen generates a vfsdata.go file that statically implements the given virtual filesystem.

vfsgen is simple and minimalistic. It features no configuration choices. You give it an input filesystem, and it generates an output .go file.

Features:

-	Outputs gofmt-compatible .go code.

-	Uses gzip compression internally (selectively, only for files that compress well).

Installation
------------

```bash
go get -u github.com/shurcooL/vfsgen
```

Usage
-----

This code will generate an assets_vfsdata.go file that statically implements the contents of "assets" directory.

```Go
var fs http.FileSystem = http.Dir("assets")

config := vfsgen.Config{
	Input: fs,
}

err := vfsgen.Generate(config)
if err != nil {
	log.Fatalln(err)
}
```

It is typically meant to be executed via go generate directives. This code can go in an assets_gen.go file, which can then be invoked via "//go:generate go run assets_gen.go". The input virtual filesystem can read directly from disk, or it can be something more involved.

By using build tags, you can create a development mode where assets are loaded directly from disk via `http.Dir`, but then statically implemented for final releases.

See [shurcooL/Go-Package-Store#38](https://github.com/shurcooL/Go-Package-Store/pull/38) for a complete example of such use.

Attribution
-----------

This package was originally based on the excellent work by [@jteeuwen](https://github.com/jteeuwen) on [`go-bindata`](https://github.com/jteeuwen/go-bindata) and [@elazarl](https://github.com/elazarl) on [`go-bindata-assetfs`](https://github.com/elazarl/go-bindata-assetfs).

License
-------

-	[MIT License](http://opensource.org/licenses/mit-license.php)
