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

var reflectPrimitiveTypes = []typed.Typeable{
	&typed.BooleanTyped{},
	&typed.NumberTyped{},
	&typed.StringTyped{},
}

func Value(val reflect.Value, targetType typed.Typeable) (typed.Valuable, error) {
	for val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	targetTypeHint := targetType.Hint()
	if typed.IsPrimitive(targetTypeHint) {
		return Primitive(val, targetType)
	}

	switch targetTypeHint {
	case typed.List:
		return List(val, targetType)
	case typed.Map:
		return Map(val, targetType)
	case typed.Object:
		return Object(val, targetType)
	case typed.Tuple:
		return Tuple(val, targetType)
	case typed.Variant:
		return Variant(val, targetType)
	default:
		return nil, fmt.Errorf("invalid type: %s", targetType.Name())
	}
}

func ValueFromTypeHint(val reflect.Value, targetTypeHint typed.TypeHint) (typed.Valuable, error) {
	for val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if typed.IsPrimitive(targetTypeHint) {
		primitveType := reflectPrimitiveTypes[targetTypeHint-1]
		return Primitive(val, primitveType)
	}

	switch targetTypeHint {
	case typed.List:
		reflectElemVal := reflect.ValueOf(val.Index(0).Interface())
		elemTypeHint := typed.TypeHintOf(reflectElemVal)
		elemVal, err := ValueFromTypeHint(reflectElemVal, elemTypeHint)
		if err == nil {
			return Value(val, &typed.ListTyped{ElementsType: elemVal.Type()})
		}
		return nil, err
	case typed.Map:
		reflectElemVal := reflect.ValueOf(val.MapIndex(val.MapKeys()[0]).Interface())
		elemTypeHint := typed.TypeHintOf(reflectElemVal)
		elemVal, err := ValueFromTypeHint(reflectElemVal, elemTypeHint)
		if err == nil {
			return Value(val, &typed.MapTyped{ElementsType: elemVal.Type()})
		}
		return nil, err
	case typed.Tuple:
		elemLen := val.Len()
		elements := make([]typed.Valuable, elemLen)
		elemTypes := make([]typed.Typeable, elemLen)
		for idx := 0; idx < elemLen; idx++ {
			reflectElemVal := reflect.ValueOf(val.Index(idx).Interface())
			elemTypeHint := typed.TypeHintOf(reflectElemVal)
			elemVal, err := ValueFromTypeHint(reflectElemVal, elemTypeHint)
			if err != nil {
				return nil, err
			}
			elements[idx] = elemVal
			elemTypes[idx] = elemVal.Type()
		}
		return &typed.TupleValue{Elements: elements, ElementsTypes: elemTypes}, nil
	case typed.Generic:
		return Generic(val)
	default:
		return nil, fmt.Errorf("invalid type hint: %v", targetTypeHint)
	}
}
