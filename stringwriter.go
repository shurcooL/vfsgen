package vfsgen

import (
	"io"
)

const lowerHex = "0123456789abcdef"

type StringWriter struct {
	io.Writer
	c int64
}

func (w *StringWriter) Write(p []byte) (n int, err error) {
	if len(p) == 0 {
		return
	}

	buf := []byte(`\x__`)
	var b byte

	for n, b = range p {
		buf[2] = lowerHex[b/16]
		buf[3] = lowerHex[b%16]
		w.Writer.Write(buf)
		w.c++
	}

	n++

	return
}
