package streams

import (
	"fmt"
	"os"
)

type OutputStream struct {
	destination *os.File
}

func (oStream *OutputStream) Write(output ...interface{}) error {
	_, err := fmt.Fprint(oStream.destination, output...)
	return err
}

func (oStream *OutputStream) Writeln(output ...interface{}) error {
	_, err := fmt.Fprintln(oStream.destination, output...)
	return err
}

func (oStream *OutputStream) Writef(format string, output ...interface{}) error {
	_, err := fmt.Fprintf(oStream.destination, format, output...)
	return err
}

func NewOutputStream(dst *os.File) *OutputStream {
	return &OutputStream{dst}
}
