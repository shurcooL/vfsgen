vfsgen
======

Package vfsgen generates a vfsdata.go file that statically implements the given virtual filesystem.

vfsgen is simple and minimalistic. It features no configuration choices. You give it an input filesystem, and it generates the output .go file.

Features:

-	Outputs gofmt-compatible .go code.
-	Supports gzip compression internally (for files that are worthwhile compressing?).

Usage
-----

TODO.

Attribution
-----------

This package was originally based on the excellent work by [@jteeuwen](https://github.com/jteeuwen) on [`go-bindata`](https://github.com/jteeuwen/go-bindata) and [@elazarl](https://github.com/elazarl) on [`go-bindata-assetfs`](https://github.com/elazarl/go-bindata-assetfs).

License
-------

- [MIT License](http://opensource.org/licenses/mit-license.php)
