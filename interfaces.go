package vfsgen

// GzipByter is implemented by compressed files for
// efficient direct access to the internal compressed bytes.
type GzipByter interface {
	// GzipBytes returns gzip compressed contents of the file.
	GzipBytes() []byte
}

// TODO: Choose one from below.

// NotWorthGzipCompressing is implemented by files that were determined
// not to be worth gzip compressing (the file size did not decrease as a result).
type NotWorthGzipCompressing interface {
	NotWorthGzipCompressing()
}

// NotWorthGzipCompressing indicates the file is not worth gzip compressing.
type NotWorthGzipCompressing_ interface {
	NotWorthGzipCompressing()
}
