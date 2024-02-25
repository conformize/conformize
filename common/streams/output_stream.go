package streams

import (
	"bytes"
	"fmt"
	"io"
)

type OutputStream struct {
	writer io.Writer
}

func (o *OutputStream) Write(p string) (int, error) {
	n, err := io.WriteString(o.writer, p)
	return n, err
}

func (o *OutputStream) Writeln(lines ...string) (int, error) {
	var buf bytes.Buffer
	buf.Grow(512)
	for _, line := range lines {
		buf.WriteString(line)
	}
	buf.WriteString("\n")
	n, err := io.WriteString(o.writer, buf.String())
	return n, err
}

func (o *OutputStream) Writef(format string, args ...any) (int, error) {
	return o.Write(fmt.Sprintf(format, args...))
}

func NewOutputStream(w io.Writer) *OutputStream {
	return &OutputStream{writer: w}
}
