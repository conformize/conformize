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

func (b *BitSet) Set(pos int) error {
	b.lock.Lock()
	defer b.lock.Unlock()
	if err := b.checkBitIsInRange(pos); err != nil {
		return err
	}
	word, bit := b.getWordAndBit(pos)
	b.bits[word] |= 1 << bit
	return nil
}

func (b *BitSet) IsSet(pos int) (bool, error) {
	b.lock.RLock()
	defer b.lock.RUnlock()
	if err := b.checkBitIsInRange(pos); err != nil {
		return false, err
	}
	word, bit := b.getWordAndBit(pos)
	return (b.bits[word] & (1 << bit)) != 0, nil
}

func (b *BitSet) Clear(pos int) error {
	b.lock.Lock()
	defer b.lock.Unlock()
	if err := b.checkBitIsInRange(pos); err != nil {
		return err
	}
	word, bit := b.getWordAndBit(pos)
	b.bits[word] &^= 1 << bit
	return nil
}

func (b *BitSet) checkBitIsInRange(bit int) error {
	if bit < 0 || bit >= b.size {
		return fmt.Errorf("index %d out of range [0, %d]", bit, b.size)
	}
	return nil
}

func (b *BitSet) getWordAndBit(pos int) (int, uint) {
	word, bit := pos/8, uint(pos%8)
	return word, bit
}
