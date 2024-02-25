package streams

import (
	"bufio"
	"os"
)

var (
	outputStream *OutputStream
	errorStream  *OutputStream
)

func Output() *OutputStream {
	return outputStream
}

func Error() *OutputStream {
	return errorStream
}

func init() {
	outputStream = NewOutputStream(bufio.NewWriterSize(os.Stdout, 2048))
	errorStream = NewOutputStream(bufio.NewWriterSize(os.Stderr, 2048))
}
