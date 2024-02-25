// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package streams

import (
	"fmt"
	"sync/atomic"

	"github.com/conformize/conformize/common/ds"
)

type Stream[T any] struct {
	closed atomic.Bool
	buf    *ds.CircularBuffer[T]
}

func (stream *Stream[T]) Read() (T, error) {
	var elem T
	elem = stream.buf.Read()
	return elem, nil
}

func (stream *Stream[T]) Write(elem T) error {
	if !stream.Closed() {
		stream.buf.Write(elem)
		return nil
	}
	return fmt.Errorf("couldn't write to stream. stream is closed")
}

func (stream *Stream[T]) Close() {
	if stream.closed.Load() {
		return
	}
	stream.closed.Store(true)
}

func (stream *Stream[T]) Closed() bool {
	return stream.closed.Load()
}

func (stream *Stream[T]) IsEmpty() bool {
	return stream.buf.IsEmpty()
}

func NewStream[T any](size int) *Stream[T] {
	diagsStream := &Stream[T]{
		buf:    ds.NewCircularBuffer[T](size),
		closed: atomic.Bool{},
	}
	return diagsStream
}
