package vfsgen

import "time"

// fileInfo is ...
type fileInfo struct {
	name             string
	uncompressedSize int64
	modTime          time.Time
}

// dirInfo is ...
type dirInfo struct {
	name    string
	entries []string
	modTime time.Time
}
