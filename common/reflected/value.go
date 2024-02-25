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

func Value(val interface{}, targetType typed.Typeable) (typed.Valuable, error) {
	targetTypeHint := targetType.Hint()
	reflectVal := reflect.ValueOf(val)
	for reflectVal.Kind() == reflect.Ptr {
		reflectVal = reflectVal.Elem()
	}
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

func ValueFromTypeHint(val interface{}, targetTypeHint typed.TypeHint) (typed.Valuable, error) {
	reflectVal := reflect.ValueOf(val)
	if reflectVal.Kind() == reflect.Ptr {
		val = reflectVal.Elem().Interface()
		return ValueFromTypeHint(val, targetTypeHint)
	}

	switch targetTypeHint {
	case typed.String:
		return Value(val, &typed.StringTyped{})
	case typed.Number:
		return Value(val, &typed.NumberTyped{})
	case typed.Boolean:
		return Value(val, &typed.BooleanTyped{})
	case typed.List:
		reflectElemVal := reflectVal.Index(0).Interface()
		elemTypeHint := typed.TypeHintOf(reflectElemVal)
		elemVal, err := ValueFromTypeHint(reflectElemVal, elemTypeHint)
		if err == nil {
			return Value(val, &typed.ListTyped{ElementsType: elemVal.Type()})
		}
		return nil, err
	case typed.Map:
		reflectElemVal := reflectVal.MapIndex(reflectVal.MapKeys()[0]).Interface()
		elemTypeHint := typed.TypeHintOf(reflectElemVal)
		elemVal, err := ValueFromTypeHint(reflectElemVal, elemTypeHint)
		if err == nil {
			return Value(val, &typed.MapTyped{ElementsType: elemVal.Type()})
		}
		return nil, err
	case typed.Tuple:
		elemLen := reflectVal.Len()
		elements := make([]typed.Valuable, elemLen)
		elemTypes := make([]typed.Typeable, elemLen)
		for idx := 0; idx < elemLen; idx++ {
			reflectElemVal := reflectVal.Index(idx).Interface()
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
		elemVal := reflectVal.Interface()
		return Generic(elemVal)
	default:
		return nil, fmt.Errorf("invalid type hint: %v", targetTypeHint)
	}
}
