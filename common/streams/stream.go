// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package streams

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/conformize/conformize/common/ds"
)

type Stream[T any] struct {
	wg         *sync.WaitGroup
	closed     atomic.Bool
	buf        *ds.CircularBuffer[T]
	streamChan chan T
	signalChan chan struct{}
}

func (stream *Stream[T]) Write(elem T) error {
	if !stream.Closed() {
		stream.buf.Write(elem)
		stream.wg.Add(1)
		stream.signalChan <- struct{}{}
		return nil
	}
	return fmt.Errorf("couldn't write to stream. stream is closed")
}

func (stream *Stream[T]) Closed() bool {
	return stream.closed.Load()
}

func (stream *Stream[T]) Close() {
	if stream.Closed() {
		return
	}
	stream.closed.Store(true)
	stream.signalChan <- struct{}{}
}

func (stream *Stream[T]) run() {
	go func() {
		defer close(stream.signalChan)
		for !stream.Closed() {
			<-stream.signalChan
			for isDrained := stream.buf.IsEmpty(); !isDrained; isDrained = stream.buf.IsEmpty() {
				stream.streamChan <- stream.buf.Read()
			}
		}
	}()
}

func NewStream[T any](wg *sync.WaitGroup, streamChan chan T, size int) *Stream[T] {
	diagsStream := &Stream[T]{
		wg:         wg,
		buf:        ds.NewCircularBuffer[T](size),
		streamChan: streamChan,
		closed:     atomic.Bool{},
		signalChan: make(chan struct{}, size),
	}
	diagsStream.run()
	return diagsStream
}
