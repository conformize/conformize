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

func (lv *ListValue) As(dst interface{}) error {
	if dst == nil {
		return fmt.Errorf("destination cannot be nil")
	}
	targetVal := reflect.ValueOf(dst)
	if targetVal.Kind() != reflect.Ptr {
		return fmt.Errorf("destination must be a pointer")
	}

	targetElem := targetVal.Elem()
	elemType := reflect.ValueOf(targetElem.Interface()).Type().Elem()
	if elemType.Kind() == reflect.Interface {
		typeHint := TypeHintFromTypeName(lv.ElementsType.Name())
		if typeHint == Invalid {
			return fmt.Errorf("no type hint defined for type %s", lv.ElementsType.Name())
		}
		resolvedNativeType, err := NativeTypeForTypeHint(typeHint)
		if err != nil {
			return err
		}
		elemType = resolvedNativeType
	}
	elements := reflect.MakeSlice(reflect.SliceOf(elemType), 0, len(lv.Elements))

	for _, element := range lv.Elements {
		elementVal := reflect.New(elemType)
		if err := element.As(elementVal.Interface()); err != nil {
			return err
		}
		elements = reflect.Append(elements, elementVal.Elem())
	}
	targetElem.Set(elements)
	return nil
}

func (lv *ListValue) String() string {
	return "ListValue"
}

func NewListValue(elements []Valuable, elementsType Typeable) Valuable {
	return &ListValue{Elements: elements, ElementsType: elementsType}
}
