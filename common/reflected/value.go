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

	if val.Kind() == reflect.Interface {
		val = reflect.ValueOf(val.Interface())
	}

	targetTypeHint := targetType.Hint()
	if typed.IsPrimitive(targetTypeHint) {
		return Primitive(val, targetType)
	}

	switch targetTypeHint.TypeHint() {
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

func ValueFromTypeHint(val reflect.Value, targetTypeHint typed.TypeHinter) (typed.Valuable, error) {
	for val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() == reflect.Interface {
		val = reflect.ValueOf(val.Interface())
	}

	if typed.IsPrimitive(targetTypeHint) {
		primitveType := reflectPrimitiveTypes[targetTypeHint.TypeHint()-1]
		return Primitive(val, primitveType)
	}

	switch targetTypeHint.TypeHint() {
	case typed.List:
		length := val.Len()
		if length == 0 {
			return &typed.ListValue{ElementsType: &typed.GenericTyped{}, Elements: []typed.Valuable{}}, nil
		}

		elements := make([]typed.Valuable, length)

		idx := 0
		elemTypeHint := typed.TypeHintOf(val.Index(idx))
		elemVal, err := ValueFromTypeHint(val.Index(idx), elemTypeHint)
		if err != nil {
			return nil, fmt.Errorf("could not reflect list element value: %v", err)
		}

		elements[idx] = elemVal
		elemType := elemVal.Type()
		idx++

		var reflectElemVal reflect.Value
		for idx < length {
			reflectElemVal = val.Index(idx)
			elemVal, err = Value(reflectElemVal, elemType)
			if err != nil {
				return nil, fmt.Errorf("could not reflect element value")
			}
			elements[idx] = elemVal
			idx++
		}
		return &typed.ListValue{Elements: elements, ElementsType: elemType}, nil
	case typed.Map:
		keys := val.MapKeys()
		elements := make(map[string]typed.Valuable, len(keys))
		var elemType typed.Typeable
		for _, key := range keys {
			if key.Kind() != reflect.String {
				continue
			}
			elemVal := val.MapIndex(key)
			elemTypeHint := typed.TypeHintOf(elemVal)
			elemTypedVal, err := ValueFromTypeHint(elemVal, elemTypeHint)
			if err != nil {
				return nil, fmt.Errorf("could not reflect map element value: %v", err)
			}
			elements[key.String()] = elemTypedVal
			if elemType == nil {
				elemType = elemTypedVal.Type()
			}
		}
		return &typed.MapValue{Elements: elements, ElementsType: elemType}, nil
	case typed.Object:
		fields := make(map[string]typed.Typeable)
		values := make(map[string]typed.Valuable)
		switch val.Kind() {
		case reflect.Map:
			for _, key := range val.MapKeys() {
				if key.Kind() != reflect.String {
					continue
				}
				fieldVal := val.MapIndex(key)
				fieldHint := typed.TypeHintOf(fieldVal)
				fieldTypedVal, err := ValueFromTypeHint(fieldVal, fieldHint)
				if err != nil {
					return nil, fmt.Errorf("could not reflect object field value: %v", err)
				}
				fields[key.String()] = fieldTypedVal.Type()
				values[key.String()] = fieldTypedVal
			}
		default:
			return nil, fmt.Errorf("cannot guess object fields from value of type %v", val.Type())
		}
		return &typed.ObjectValue{Fields: values, FieldsTypes: fields}, nil

	case typed.Tuple:
		length := val.Len()
		elements := make([]typed.Valuable, length)
		elemTypes := make([]typed.Typeable, length)
		for i := 0; i < length; i++ {
			elemVal := val.Index(i)
			elemTypeHint := typed.TypeHintOf(elemVal)
			elemTypedVal, err := ValueFromTypeHint(elemVal, elemTypeHint)
			if err != nil {
				return nil, fmt.Errorf("could not reflect tuple element value: %v", err)
			}
			elements[i] = elemTypedVal
			elemTypes[i] = elemTypedVal.Type()
		}
		return &typed.TupleValue{Elements: elements, ElementsTypes: elemTypes}, nil

	case typed.Generic:
		return Generic(val)
	default:
		return nil, fmt.Errorf("invalid type hint: %v", targetTypeHint)
	}
}
