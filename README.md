# vfsgen [![Build Status](https://travis-ci.org/shurcooL/vfsgen.svg?branch=master)](https://travis-ci.org/shurcooL/vfsgen) [![GoDoc](https://godoc.org/github.com/shurcooL/vfsgen?status.svg)](https://godoc.org/github.com/shurcooL/vfsgen)

Package vfsgen generates a vfsdata.go file that statically implements the given virtual filesystem.

vfsgen is simple and minimalistic. It features no configuration choices. You give it an input filesystem, and it generates an output .go file.

Features:

-	Outputs gofmt-compatible .go code.

-	Supports gzip compression internally.

Installation
------------

```bash
go get -u github.com/shurcooL/vfsgen
```

Usage
-----

vfsgen is great to use via go generate directives. By using build tags, you can create a development mode where assets are loaded directly from disk via `http.Dir`, but then statically implemented for final releases.

See [shurcooL/Go-Package-Store#38](https://github.com/shurcooL/Go-Package-Store/pull/38) for a complete example of such use.

Attribution
-----------

This package was originally based on the excellent work by [@jteeuwen](https://github.com/jteeuwen) on [`go-bindata`](https://github.com/jteeuwen/go-bindata) and [@elazarl](https://github.com/elazarl) on [`go-bindata-assetfs`](https://github.com/elazarl/go-bindata-assetfs).

License
-------

- [MIT License](http://opensource.org/licenses/mit-license.php)
