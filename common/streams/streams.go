package streams

import (
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
	outputStream = NewOutputStream(os.Stdout)
	errorStream = NewOutputStream(os.Stderr)
}
