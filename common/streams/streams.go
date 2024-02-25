package streams

import (
	"os"
	"sync"
)

type streams struct {
	stdOut *OutputStream
	stdErr *OutputStream
}

func (s *streams) Output() *OutputStream {
	return s.stdOut
}

func (s *streams) Error() *OutputStream {
	return s.stdErr
}

var (
	instance *streams
	once     sync.Once
)

func Instance() *streams {
	once.Do(func() {
		instance = &streams{
			stdOut: NewOutputStream(os.Stdout),
			stdErr: NewOutputStream(os.Stderr),
		}
	})
	return instance
}
