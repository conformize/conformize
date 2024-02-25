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
	"unsafe"
)

type CircularBuffer[T any] struct {
	buffer              []unsafe.Pointer
	size                int
	start               atomic.Int64
	end                 atomic.Int64
	signalNotEmpty      *sync.Cond
	signalAvailableSlot *sync.Cond
}

func NewCircularBuffer[T any](size int) *CircularBuffer[T] {
	cb := &CircularBuffer[T]{
		buffer:              make([]unsafe.Pointer, size),
		size:                size,
		start:               atomic.Int64{},
		end:                 atomic.Int64{},
		signalNotEmpty:      sync.NewCond(&sync.Mutex{}),
		signalAvailableSlot: sync.NewCond(&sync.Mutex{}),
	}
	cb.start.Store(0)
	cb.end.Store(0)

	return cb
}

func (cb *CircularBuffer[T]) Write(elem T) {
	cb.signalAvailableSlot.L.Lock()
	for cb.IsFull() {
		cb.signalAvailableSlot.Wait()
	}
	cb.signalAvailableSlot.L.Unlock()

	index := cb.end.Load() % int64(cb.size)
	cb.end.Add(1)
	atomic.StorePointer(&cb.buffer[index], unsafe.Pointer(&elem))
	cb.signalNotEmpty.L.Lock()
	cb.signalNotEmpty.Signal()
	cb.signalNotEmpty.L.Unlock()
}

func (cb *CircularBuffer[T]) Read() T {
	cb.signalNotEmpty.L.Lock()
	for cb.IsEmpty() {
		cb.signalNotEmpty.Wait()
	}
	cb.signalNotEmpty.L.Unlock()

	var elem T
	index := cb.start.Load() % int64(cb.size)
	cb.start.Add(1)
	elem = *(*T)(atomic.LoadPointer(&cb.buffer[index]))
	cb.signalAvailableSlot.L.Lock()
	cb.signalAvailableSlot.Signal()
	cb.signalAvailableSlot.L.Unlock()
	return elem
}

func (cb *CircularBuffer[T]) IsEmpty() bool {
	return cb.end.Load() == cb.start.Load()
}

func (cb *CircularBuffer[T]) IsFull() bool {
	return (cb.end.Load() - cb.start.Load()) >= int64(cb.size)
}
