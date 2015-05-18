/*
Package vfsgen generates a vfsdata.go file that statically implements the given virtual filesystem.

vfsgen is simple and minimalistic. It features no configuration choices. You give it an input filesystem, and it generates the output .go file.

Features:

-	Outputs gofmt-compatible .go code.

-	Supports gzip compression internally.
*/
package vfsgen
