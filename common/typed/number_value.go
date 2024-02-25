// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package typed

import "fmt"

type NumberValue struct {
	PrimitiveValue[float64]
}

func (numVal *NumberValue) Type() Typeable {
	return &NumberTyped{}
}

func (numVal *NumberValue) Assign(val Valuable) error {
	intVal, ok := val.(*NumberValue)
	if !ok {
		return fmt.Errorf("cannot apply %v to %v", val, numVal)
	}
	*numVal = *intVal
	return nil
}

func (numVal *NumberValue) String() string {
	return "NumberValue"
}

func NewNumberValue(value interface{}) (Valuable, error) {
	v, err := NewPrimitive[float64](value)
	return &NumberValue{*v}, err
}
