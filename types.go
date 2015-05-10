package vfsgen

import "time"

// file is ...
type file struct {
	name             string
	uncompressedSize int64
	modTime          time.Time
}

// dir is ...
type dir struct {
	name    string
	entries []string
	modTime time.Time
}
