// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package ds

import (
	"sync"
	"sync/atomic"
)

type CircularBuffer[T any] struct {
	buffer              []T
	size                int
	start               atomic.Int32
	end                 atomic.Int32
	count               atomic.Int32
	signalNotEmpty      *sync.Cond
	signalAvailableSlot *sync.Cond
}

func NewCircularBuffer[T any](size int) *CircularBuffer[T] {
	cb := &CircularBuffer[T]{
		buffer:              make([]T, size),
		size:                size,
		start:               atomic.Int32{},
		end:                 atomic.Int32{},
		count:               atomic.Int32{},
		signalNotEmpty:      sync.NewCond(&sync.Mutex{}),
		signalAvailableSlot: sync.NewCond(&sync.Mutex{}),
	}
	return cb
}

func (cb *CircularBuffer[T]) Write(elem T) {
	cb.signalAvailableSlot.L.Lock()
	for int(cb.count.Load()) == cb.size {
		cb.signalAvailableSlot.Wait()
	}
	cb.signalAvailableSlot.L.Unlock()

	end := cb.end.Load()
	cb.buffer[end] = elem

	cb.end.CompareAndSwap(end, (end+1)%int32(cb.size))
	cb.count.Add(1)
	cb.signalNotEmpty.Signal()
}

func (cb *CircularBuffer[T]) Read() T {
	cb.signalNotEmpty.L.Lock()
	for cb.count.Load() == 0 {
		cb.signalNotEmpty.Wait()
	}
	cb.signalNotEmpty.L.Unlock()

	start := cb.start.Load()
	diag := cb.buffer[start]

	start = cb.start.Load()
	cb.start.CompareAndSwap(start, (start+1)%int32(cb.size))
	cb.count.Add(-1)
	cb.signalAvailableSlot.Signal()
	return diag
}

func (cb *CircularBuffer[T]) IsEmpty() bool {
	return cb.count.Load() == 0
}
