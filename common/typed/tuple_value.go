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

type TupleValue struct {
	Elements      []Valuable
	ElementsTypes []Typeable
}

func (tv *TupleValue) Type() Typeable {
	return &TupleTyped{ElementsTypes: tv.ElementsTypes}
}

func (tv *TupleValue) Assign(val Valuable) error {
	tupVal, ok := val.(*TupleValue)
	if !ok {
		return fmt.Errorf("cannot apply %v to %v", val, tv)
	}
	*tv = *tupVal
	return nil
}

func (tv *TupleValue) As(dst interface{}) error {
	if dst == nil {
		return fmt.Errorf("destination cannot be nil")
	}
	targetVal := reflect.ValueOf(dst)
	if targetVal.Kind() != reflect.Ptr {
		return fmt.Errorf("destination must be a pointer")
	}

	sliceVal := reflect.ValueOf(targetVal.Elem().Interface())
	sliceLen := sliceVal.Len()
	if sliceLen > 0 && sliceLen != len(tv.Elements) {
		return fmt.Errorf("cannot reflect tuple to slice - mismatching number of elements")
	}

	var elemType reflect.Type
	elemType = reflect.TypeOf((*interface{})(nil)).Elem()
	elements := reflect.MakeSlice(reflect.SliceOf(elemType), 0, len(tv.Elements))

	useNativeTypeForTypeHint := sliceLen == 0
	for idx, element := range tv.Elements {
		var elemVal reflect.Value
		if useNativeTypeForTypeHint {
			elemNativeType, err := NativeTypeForTypeHint(element.Type().Hint())
			if err != nil {
				return err
			}
			elemType = elemNativeType
		} else {
			elemType = reflect.ValueOf(sliceVal.Index(idx).Interface()).Type()
		}
		elemVal = reflect.New(elemType)
		if err := element.As(elemVal.Interface()); err != nil {
			return err
		}
		elements = reflect.Append(elements, elemVal.Elem())
	}
	targetVal.Elem().Set(elements)
	return nil
}

func (tv *TupleValue) Length() int {
	return len(tv.Elements)
}

func (tv *TupleValue) String() string {
	return "TupleValue"
}
