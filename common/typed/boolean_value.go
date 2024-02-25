// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package typed

import (
	"fmt"
	"reflect"
)

var _booleanValueType = &BooleanTyped{}

type BooleanValue struct {
	PrimitiveValue[bool]
}

func (bv *BooleanValue) Type() Typeable {
	return _booleanValueType
}

func (bv *BooleanValue) Assign(val Valuable) error {
	boolVal, ok := val.(*BooleanValue)
	if !ok {
		return fmt.Errorf("cannot apply %v to %v", val, bv)
	}
	bv.PrimitiveValue = boolVal.PrimitiveValue
	return nil
}

func (bv *BooleanValue) String() string {
	return "BooleanValue"
}

func NewBooleanValue(value any) (Valuable, error) {
	v, err := NewPrimitive[bool](reflect.ValueOf(value))
	return &BooleanValue{*v}, err
}
