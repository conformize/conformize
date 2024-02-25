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

func (tv *TupleValue) As(dst any) error {
	if dst == nil {
		return fmt.Errorf("destination cannot be nil")
	}

	dstVal := reflect.ValueOf(dst)
	if dstVal.Kind() != reflect.Ptr || dstVal.IsNil() {
		return fmt.Errorf("destination must be a non-nil pointer")
	}

	sliceVal := dstVal.Elem()
	if sliceVal.Kind() == reflect.Interface {
		if sliceVal.IsNil() {
			return fmt.Errorf("destination interface is nil")
		}
		sliceVal = reflect.ValueOf(sliceVal.Interface())
		if sliceVal.Kind() != reflect.Slice {
			return fmt.Errorf("interface must hold a slice, got %s", sliceVal.Kind())
		}
	}

	sliceLen := sliceVal.Len()
	if sliceLen > 0 && sliceLen != len(tv.Elements) {
		return fmt.Errorf("cannot reflect tuple to slice - mismatching number of elements")
	}

	var elemType reflect.Type
	elemType = reflect.TypeOf((*any)(nil)).Elem()
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
			elemType = sliceVal.Index(idx).Type()
			if elemType.Kind() == reflect.Interface {
				if sliceVal.Index(idx).IsNil() {
					return fmt.Errorf("element at index %d is nil", idx)
				}

				elemType = reflect.ValueOf(sliceVal.Index(idx).Interface()).Type()
			}
		}

		elemVal = reflect.New(elemType)
		if err := element.As(elemVal.Interface()); err != nil {
			return err
		}

		elements = reflect.Append(elements, elemVal.Elem())
	}

	if elements.CanConvert(sliceVal.Type().Elem()) {
		elements = elements.Convert(sliceVal.Type().Elem())
	}

	dstVal.Elem().Set(elements)
	return nil
}

func (tv *TupleValue) Length() int {
	return len(tv.Elements)
}

func (tv *TupleValue) String() string {
	return "TupleValue"
}

func (tv *TupleValue) Items() []Valuable {
	return tv.Elements
}
