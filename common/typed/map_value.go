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

type MapValue struct {
	Elements     map[string]Valuable
	ElementsType Typeable
}

func (mv *MapValue) Type() Typeable {
	return &MapTyped{ElementsType: mv.ElementsType}
}

func (mv *MapValue) Assign(val Valuable) error {
	mapVal, ok := val.(*MapValue)
	if !ok {
		return fmt.Errorf("cannot apply %v to %v", val, mv)
	}
	*mv = *mapVal
	return nil
}

func (mv *MapValue) As(dst any) error {
	if dst == nil {
		return fmt.Errorf("destination cannot be nil")
	}
	dstVal := reflect.ValueOf(dst)

	mapVal := reflect.MakeMap(dstVal.Type().Elem())
	elemType := dstVal.Type().Elem().Elem()
	for k, v := range mv.Elements {
		vVal := reflect.New(elemType)
		err := v.As(vVal.Interface())
		if err != nil {
			return err
		}
		mapVal.SetMapIndex(reflect.ValueOf(k), vVal.Elem())
	}
	dstVal.Elem().Set(mapVal)
	return nil
}

func (mv *MapValue) String() string {
	return "MapValue"
}

func (mv *MapValue) Lenght() int {
	return len(mv.Elements)
}

func NewMapValue(elements map[string]Valuable, elementsType Typeable) Valuable {
	return &MapValue{Elements: elements, ElementsType: elementsType}
}
