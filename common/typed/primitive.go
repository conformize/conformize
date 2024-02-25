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

type primitiveType interface {
	bool | float64 | string
}

type PrimitiveValue[T primitiveType] struct {
	value T
}

func NewPrimitive[T primitiveType](value reflect.Value) (*PrimitiveValue[T], error) {
	if !value.IsValid() || (value.Kind() == reflect.Ptr && value.IsNil()) {
		return &PrimitiveValue[T]{}, nil
	}

	var val T
	reflectVal := reflect.ValueOf(val)
	if reflectVal.Kind() != value.Kind() {
		if !value.Type().ConvertibleTo(reflectVal.Type()) {
			return nil, fmt.Errorf("cannot convert value of type %s to type %s",
				value.Type().String(), reflectVal.Type().String(),
			)
		}
		reflectVal = value.Convert(reflectVal.Type())
	} else {
		reflectVal = value
	}

	var ok bool
	val, ok = reflect.TypeAssert[T](reflectVal)
	if !ok {
		return nil, fmt.Errorf("cannot use value of type %s as type %T",
			reflectVal.Type().String(), val,
		)
	}
	return &PrimitiveValue[T]{value: val}, nil
}

func (pv *PrimitiveValue[T]) As(dst any) error {
	if dst == nil {
		return fmt.Errorf("destination cannot be nil")
	}

	dstVal := reflect.ValueOf(dst)
	if dstVal.Kind() != reflect.Ptr || dstVal.IsNil() {
		return fmt.Errorf("destination must be a non-nil pointer")
	}

	target := dstVal.Elem()
	src := reflect.ValueOf(pv.value)
	if target.Kind() == reflect.Interface && !target.IsNil() {
		target = reflect.ValueOf(target.Interface())
	}

	targetType := target.Type()
	if !src.Type().ConvertibleTo(targetType) {
		return fmt.Errorf("cannot convert %s to %s", src.Type(), targetType)
	}

	dstVal.Elem().Set(src.Convert(targetType))
	return nil
}
