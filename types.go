package vfsgen

import "time"

// compressedFile is ...
type compressedFile struct {
	name              string
	compressedContent []byte
	uncompressedSize  int64
	modTime           time.Time
}

// dir is ...
type dir struct {
	name string
	//entries []os.FileInfo
	entries []string
	modTime time.Time
}
