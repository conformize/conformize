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
	return &PrimitiveValue[T]{value: reflectVal.Interface().(T)}, nil
}

func (pv *PrimitiveValue[T]) As(dst interface{}) error {
	if dst == nil {
		return fmt.Errorf("destination cannot be nil")
	}

	dstVal := reflect.ValueOf(dst)
	if dstVal.Kind() != reflect.Ptr {
		return fmt.Errorf("destination must be a pointer")
	}

	v := reflect.ValueOf(pv.value)
	valType := v.Type()
	dstType := reflect.ValueOf(dstVal.Elem().Interface()).Type()
	if valType.Kind() != dstType.Kind() {
		if !valType.ConvertibleTo(dstType) {
			return fmt.Errorf("cannot convert value of type %s to type %s", valType.String(), dstType.String())
		}
		v = v.Convert(dstType)
	}
	dstVal.Elem().Set(v)
	return nil
}
