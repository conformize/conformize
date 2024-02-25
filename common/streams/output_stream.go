package streams

import (
	"bufio"
	"fmt"
	"io"
)

type OutputStream struct {
	writer io.Writer
	flush  func() error
}

func (o *OutputStream) Write(p string) (int, error) {
	n, err := io.WriteString(o.writer, p)
	o.Flush()
	return n, err
}

func (o *OutputStream) Writeln(lines ...string) (int, error) {
	for _, line := range lines {
		io.WriteString(o.writer, line)
	}
	n, err := io.WriteString(o.writer, "\n")
	o.Flush()
	return n, err
}

func (o *OutputStream) Writef(format string, args ...interface{}) (int, error) {
	return o.Write(fmt.Sprintf(format, args...))
}

func (o *OutputStream) Flush() error {
	return o.flush()
}

func NewOutputStream(w io.Writer) *OutputStream {
	out := &OutputStream{
		writer: w,
		flush:  func() error { return nil },
	}
	if bw, ok := w.(*bufio.Writer); ok {
		out.flush = bw.Flush
	}
	return out
}
