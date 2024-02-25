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

type TypeHint byte

const (
	Boolean TypeHint = 0b00000001
	Number  TypeHint = 0b00000010
	String  TypeHint = 0b00000011
	List    TypeHint = 0b10000000
	Map     TypeHint = 0b11000000
	Object  TypeHint = 0b11100000
	Tuple   TypeHint = 0b11110000
	Variant TypeHint = 0b11111000
	Generic TypeHint = 0b11111100
	Invalid TypeHint = 0b11111111
)

func TypeHintFromTypeName(typeName NamedType) TypeHint {
	switch typeName {
	case BooleanType:
		return Boolean
	case NumberType:
		return Number
	case StringType:
		return String
	case ListType:
		return List
	case MapType:
		return Map
	case ObjectType:
		return Object
	case TupleType:
		return Tuple
	case VariantType:
		return Variant
	case GenericType:
		return Generic
	default:
		return Invalid
	}
}

func IsPrimitive(hint TypeHint) bool {
	return (hint>>7)&0b1 == 0
}

var typeHintMappings = map[reflect.Kind]TypeHint{
	reflect.Bool:      Boolean,
	reflect.Int:       Number,
	reflect.Int8:      Number,
	reflect.Int16:     Number,
	reflect.Int32:     Number,
	reflect.Int64:     Number,
	reflect.Uint:      Number,
	reflect.Uint8:     Number,
	reflect.Uint16:    Number,
	reflect.Uint32:    Number,
	reflect.Uint64:    Number,
	reflect.Float32:   Number,
	reflect.Float64:   Number,
	reflect.String:    String,
	reflect.Array:     List,
	reflect.Slice:     List,
	reflect.Map:       Map,
	reflect.Struct:    Object,
	reflect.Interface: Generic,
}

func TypeHintOf(v interface{}) TypeHint {
	reflectVal := reflect.ValueOf(v)
	valueKind := reflectVal.Type().Kind()
	if isTuple(reflectVal.Interface()) {
		return Tuple
	}
	return typeHintMappings[valueKind]
}

func isTuple(val interface{}) bool {
	v := reflect.ValueOf(val)
	if v.Kind() != reflect.Slice {
		return false
	}

	if v.Len() == 0 {
		return true
	}

	elemType := reflect.TypeOf(v.Index(0).Interface())
	sliceLen := v.Len()
	for startIdx, endIdx := 0, sliceLen-1; endIdx >= startIdx; startIdx, endIdx = startIdx+1, endIdx-1 {
		for _, idx := range []int{startIdx, endIdx} {
			if reflect.TypeOf(v.Index(idx).Interface()) != elemType {
				return true
			}
		}
		if startIdx == endIdx {
			break
		}
	}
	return false
}

func NativeTypeForTypeHint(typeHint TypeHint) (reflect.Type, error) {
	switch typeHint {
	case Boolean:
		return reflect.TypeOf(true), nil
	case Number:
		return reflect.TypeOf(float64(0)), nil
	case String:
		return reflect.TypeOf(""), nil
	case List, Tuple:
		return reflect.TypeOf([]interface{}{}), nil
	case Map:
		return reflect.TypeOf(map[interface{}]interface{}{}), nil
	case Object:
		return reflect.TypeOf(struct{}{}), nil
	case Variant, Generic:
		return reflect.TypeOf((*interface{})(nil)).Elem(), nil
	default:
		return nil, fmt.Errorf("cannot resolve type for unknown type hint %d", typeHint)
	}
}
