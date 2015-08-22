/*
Package vfsgen generates a vfsdata.go file that statically implements the given virtual filesystem.

vfsgen is simple and minimalistic. It features no configuration choices. You give it an input filesystem, and it generates an output .go file.

Features:

-	Outputs gofmt-compatible .go code.

-	Uses gzip compression internally (selectively, only for files that compress well).
*/
package vfsgen
