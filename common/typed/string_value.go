// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package typed

import "fmt"

type StringValue struct {
	PrimitiveValue[string]
}

func (stv *StringValue) Type() Typeable {
	return &StringTyped{}
}

func (stv *StringValue) Assign(val Valuable) error {
	strValue, ok := val.(*StringValue)
	if !ok {
		return fmt.Errorf("cannot apply %v to %v", val, stv)
	}
	*stv = *strValue
	return nil
}

func (stv *StringValue) String() string {
	return "StringValue"
}

func (stv *StringValue) Length() int {
	return len(stv.value)
}

func NewStringValue(value interface{}) (Valuable, error) {
	v, err := NewPrimitive[string](value)
	return &StringValue{*v}, err
}