// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package ds

import (
	"fmt"
	"sync"
)

type BitSet struct {
	bits []byte
	size int
	lock sync.RWMutex
}

func NewBitSet(size int) *BitSet {
	words := (size + 7) / 8
	return &BitSet{
		bits: make([]byte, words),
		size: size,
	}
}

func (b *BitSet) Set(index int) error {
	if index < 0 || index >= b.size {
		return fmt.Errorf("index %d out of range [0, %d)", index, b.size)
	}

	b.lock.Lock()
	defer b.lock.Unlock()
	word, bit := index/8, uint(index%8)
	b.bits[word] |= 1 << bit
	return nil
}

func (b *BitSet) IsSet(index int) (bool, error) {
	if index < 0 || index >= b.size {
		return false, fmt.Errorf("index %d out of range [0, %d)", index, b.size)

	}
	b.lock.RLock()
	defer b.lock.RUnlock()
	word, bit := index/8, uint(index%8)
	return (b.bits[word] & (1 << bit)) != 0, nil
}

func (b *BitSet) Clear(index int) error {
	if index < 0 || index >= b.size {
		return fmt.Errorf("index %d out of range [0, %d)", index, b.size)
	}
	b.lock.Lock()
	defer b.lock.Unlock()
	word, bit := index/8, uint(index%8)
	b.bits[word] &^= 1 << bit
	return nil
}
