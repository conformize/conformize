// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package typed

import (
	"fmt"
)

type BooleanValue struct {
	PrimitiveValue[bool]
}

func (bv *BooleanValue) Type() Typeable {
	return &BooleanTyped{}
}

func (bv *BooleanValue) Assign(val Valuable) error {
	boolVal, ok := val.(*BooleanValue)
	if !ok {
		return fmt.Errorf("cannot apply %v to %v", val, bv)
	}
	*bv = *boolVal
	return nil
}

func (bv *BooleanValue) String() string {
	return "BooleanValue"
}

func NewBooleanValue(value interface{}) (Valuable, error) {
	v, err := NewPrimitive[bool](value)
	return &BooleanValue{*v}, err
}
