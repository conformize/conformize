// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package reflected

import (
	"fmt"
	"reflect"

	"github.com/conformize/conformize/common/typed"
)

type primitiveFn func(reflect.Value) (typed.Valuable, error)

var primitiveFns = []primitiveFn{
	func(v reflect.Value) (typed.Valuable, error) {
		primVal, err := typed.NewPrimitive[bool](v)
		return &typed.BooleanValue{PrimitiveValue: *primVal}, err
	},
	func(v reflect.Value) (typed.Valuable, error) {
		primVal, err := typed.NewPrimitive[float64](v)
		return &typed.NumberValue{PrimitiveValue: *primVal}, err
	},
	func(v reflect.Value) (typed.Valuable, error) {
		primVal, err := typed.NewPrimitive[string](v)
		return &typed.StringValue{PrimitiveValue: *primVal}, err
	},
}

func Primitive(value reflect.Value, targetType typed.Typeable) (typed.Valuable, error) {
	if value.IsValid() {
		valueTypeHint := typed.TypeHintOf(value)
		isPrimitive := typed.IsPrimitive(targetType.Hint())
		if !isPrimitive || valueTypeHint.TypeHint() != targetType.Hint().TypeHint() {
			return nil, fmt.Errorf("invalid primitive type, expected %s", targetType.Name())
		}
	}
	primitiveFn := primitiveFns[targetType.Hint().TypeHint()-1]
	return primitiveFn(value)
}
