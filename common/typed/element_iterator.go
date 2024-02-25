// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package typed

type ElementIterator[T Elementable] struct {
	value T
	pos   int
}

func NewElementIterator[T Elementable](value T) *ElementIterator[T] {
	return &ElementIterator[T]{
		value: value,
		pos:   0,
	}
}

func (it *ElementIterator[T]) Next() bool {
	if it.pos < len(it.value.Items()) {
		it.pos++
		return true
	}
	return false
}

func (it *ElementIterator[T]) Element(fn func(v Valuable)) {
	if it.pos > 0 && it.pos <= len(it.value.Items()) {
		fn(it.value.Items()[it.pos-1])
	}
}
