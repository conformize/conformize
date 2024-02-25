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

type TypeHinter interface {
	TypeHint() TypeHint
}

type ComplexTypeHinter interface {
	TypeHinter
	ElementsTypeHint() TypeHinter
}

type ComplexTypeFieldsHinter interface {
	TypeHinter
	FieldsTypeHint() map[string]TypeHinter
}

type ComplexTypeMixedElementHinter interface {
	TypeHinter
	ElementsTypeHints() []TypeHinter
}

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

type simpleTypeHint struct {
	kind TypeHint
}

func (sth *simpleTypeHint) TypeHint() TypeHint {
	return sth.kind
}

type complexTypeHint struct {
	kind     TypeHint
	elemType TypeHinter
}

type complexObjectHint struct {
	kind       TypeHint
	fieldsType map[string]TypeHinter
}

type complexTypeMixedElementsHint struct {
	kind         TypeHint
	elementsType []TypeHinter
}

func (coh *complexObjectHint) TypeHint() TypeHint {
	return coh.kind
}

func (coh *complexObjectHint) FieldsTypeHint() map[string]TypeHinter {
	return coh.fieldsType
}

func (cth *complexTypeHint) TypeHint() TypeHint {
	return cth.kind
}

func (cth *complexTypeHint) ElementsTypeHint() TypeHinter {
	return cth.elemType
}

func (cth *complexTypeMixedElementsHint) TypeHint() TypeHint {
	return cth.kind
}

func (cth *complexTypeMixedElementsHint) ElementsTypeHints() []TypeHinter {
	return cth.elementsType
}

func (th TypeHint) String() NamedType {
	switch th {
	case Boolean:
		return BooleanType
	case Number:
		return NumberType
	case String:
		return StringType
	case List:
		return ListType
	case Map:
		return MapType
	case Object:
		return ObjectType
	case Tuple:
		return TupleType
	case Variant:
		return VariantType
	case Generic:
		return GenericType
	default:
		return InvalidType
	}
}

func IsPrimitive(hint TypeHinter) bool {
	return (hint.TypeHint()>>7)&0b1 == 0
}

var primitiveTypeHintMappings = map[reflect.Kind]TypeHint{
	reflect.Bool:    Boolean,
	reflect.Int:     Number,
	reflect.Int8:    Number,
	reflect.Int16:   Number,
	reflect.Int32:   Number,
	reflect.Int64:   Number,
	reflect.Uint:    Number,
	reflect.Uint8:   Number,
	reflect.Uint16:  Number,
	reflect.Uint32:  Number,
	reflect.Uint64:  Number,
	reflect.Float32: Number,
	reflect.Float64: Number,
	reflect.String:  String,
}

func TypeHintOf(v reflect.Value) TypeHinter {
	val := v

	for val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return &simpleTypeHint{kind: Invalid}
		}
		val = val.Elem()
	}

	if val.Kind() == reflect.Interface {
		if val.IsNil() {
			return &simpleTypeHint{kind: Invalid}
		}
		val = reflect.ValueOf(val.Interface())
	}

	kind, ok := primitiveTypeHintMappings[val.Kind()]
	if ok {
		return &simpleTypeHint{kind: kind}
	}

	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		if val.Len() == 0 {
			return &complexTypeHint{kind: List, elemType: &simpleTypeHint{kind: Generic}}
		}

		startIdx, endIdx := 0, val.Len()-1

		var elem reflect.Value
		var elemHint TypeHinter

		firstHint := TypeHintOf(val.Index(startIdx))
		typeHints := make([]TypeHinter, val.Len())
		typeHints[startIdx] = firstHint

		startIdx++
		isHeterogeneous := false
		for startIdx <= endIdx {
			elem = val.Index(startIdx)
			elemHint = TypeHintOf(elem)
			if elemHint.TypeHint() != firstHint.TypeHint() {
				isHeterogeneous = true
			}

			typeHints[startIdx] = elemHint
			if startIdx == endIdx {
				break
			}

			elem = val.Index(endIdx)
			elemHint = TypeHintOf(elem)
			if elemHint.TypeHint() != firstHint.TypeHint() {
				isHeterogeneous = true
			}
			typeHints[endIdx] = elemHint

			startIdx++
			endIdx--
		}

		if !isHeterogeneous {
			return &complexTypeHint{kind: List, elemType: firstHint}
		}
		return &complexTypeMixedElementsHint{kind: Tuple, elementsType: typeHints}

	case reflect.Map:
		if val.Len() == 0 {
			return &complexTypeHint{kind: Map, elemType: &simpleTypeHint{kind: Generic}}
		}

		elemVal := val.MapIndex(val.MapKeys()[0])
		elemHint := TypeHintOf(elemVal)
		return &complexTypeHint{kind: Map, elemType: elemHint}

	case reflect.Struct:
		fieldsType := make(map[string]TypeHinter)
		for i := 0; i < val.NumField(); i++ {
			field := val.Type().Field(i)
			fieldVal := val.Field(i)
			fieldHint := TypeHintOf(fieldVal)
			fieldsType[field.Name] = fieldHint
		}
		return &complexObjectHint{kind: Object, fieldsType: fieldsType}

	case reflect.Interface:
		return &simpleTypeHint{kind: Generic}
	default:
		return &simpleTypeHint{kind: Invalid}
	}

}

var _nativeBoolType = reflect.TypeOf(true)
var _nativeNumberType = reflect.TypeOf(float64(0))
var _nativeStringType = reflect.TypeOf("")
var _nativeGenericType = reflect.TypeOf((*any)(nil)).Elem()
var _nativeObjectType = reflect.TypeOf(struct{}{})

func NativeTypeForTypeHint(hinter TypeHinter) (reflect.Type, error) {
	switch hint := hinter.(type) {

	case *simpleTypeHint:
		switch hint.kind {
		case Boolean:
			return _nativeBoolType, nil
		case Number:
			return _nativeNumberType, nil
		case String:
			return _nativeStringType, nil
		case Object:
			return _nativeObjectType, nil
		case Variant, Generic:
			return _nativeGenericType, nil
		case Invalid:
			return nil, fmt.Errorf("cannot resolve type for invalid type hint")
		default:
			return nil, fmt.Errorf("unknown simple type hint: %v", hint.kind)
		}

	case *complexTypeHint:
		elemType, err := NativeTypeForTypeHint(hint.elemType)
		if err != nil {
			return nil, fmt.Errorf("could not resolve element type: %w", err)
		}

		switch hint.kind {
		case List, Tuple:
			return reflect.SliceOf(elemType), nil
		case Map:
			return reflect.MapOf(_nativeStringType, elemType), nil
		case Variant:
			return _nativeGenericType, nil
		default:
			return nil, fmt.Errorf("unknown complex type hint: %v", hint.kind)
		}

	default:
		return nil, fmt.Errorf("unrecognized TypeHinter implementation")
	}
}
