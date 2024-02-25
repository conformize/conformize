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

type ListValue struct {
	ElementsType Typeable
	Elements     []Valuable
}

func (lv *ListValue) Type() Typeable {
	return &ListTyped{ElementsType: lv.ElementsType}
}

func (lv *ListValue) Assign(val Valuable) error {
	listVal, ok := val.(*ListValue)
	if !ok {
		return fmt.Errorf("cannot apply %v to %v", val, lv)
	}
	*lv = *listVal
	return nil
}

func (lv *ListValue) As(dst any) error {
	if dst == nil {
		return fmt.Errorf("destination cannot be nil")
	}

	dstVal := reflect.ValueOf(dst)
	if dstVal.Kind() != reflect.Ptr || dstVal.IsNil() {
		return fmt.Errorf("destination must be a non-nil pointer")
	}

	dstElem := dstVal.Elem()
	var dstElemType reflect.Type
	switch dstElem.Kind() {
	case reflect.Slice:
		dstElemType = dstElem.Type().Elem()
		dstVal = dstElem
	case reflect.Interface:
		if dstElem.IsNil() {
			return fmt.Errorf("destination interface is nil")
		}
		concrete := reflect.ValueOf(dstElem.Interface())
		if concrete.Kind() != reflect.Slice {
			return fmt.Errorf("interface must hold a slice, got %s", concrete.Kind())
		}
		dstElemType = concrete.Type().Elem()
		dstVal = dstElem

	default:
		return fmt.Errorf("destination must be a pointer to a slice or interface holding a slice, got %s", dstElem.Kind())
	}

	if dstElemType.Kind() == reflect.Interface {
		typeHint := lv.ElementsType.Hint()
		if typeHint.TypeHint() == Invalid {
			return fmt.Errorf("no type hint defined for element type %s", lv.ElementsType.Name())
		}

		var err error
		dstElemType, err = NativeTypeForTypeHint(typeHint)
		if err != nil {
			return err
		}
	}

	newSlice := reflect.MakeSlice(reflect.SliceOf(dstElemType), 0, len(lv.Elements))
	for _, el := range lv.Elements {
		ptr := reflect.New(dstElemType)
		if err := el.As(ptr.Interface()); err != nil {
			return fmt.Errorf("failed to convert list element: %w", err)
		}
		newSlice = reflect.Append(newSlice, ptr.Elem())
	}

	dstVal.Set(newSlice)
	return nil
}

func (lv *ListValue) Items() []Valuable {
	return lv.Elements
}

func (lv *ListValue) String() string {
	return "ListValue"
}

func (lv *ListValue) Length() int {
	return len(lv.Elements)
}

func NewListValue(elements []Valuable, elementsType Typeable) Valuable {
	return &ListValue{Elements: elements, ElementsType: elementsType}
}
