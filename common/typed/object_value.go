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

type ObjectValue struct {
	Fields      map[string]Valuable
	FieldsTypes map[string]Typeable
}

func (ov *ObjectValue) Type() Typeable {
	return &ObjectTyped{FieldsTypes: ov.FieldsTypes}
}

func (ov *ObjectValue) Assign(val Valuable) error {
	objVal, ok := val.(*ObjectValue)
	if !ok {
		return fmt.Errorf("cannot apply %v to %v", val, ov)
	}
	*ov = *objVal
	return nil
}

func (ov *ObjectValue) As(dst any) error {
	if dst == nil {
		return fmt.Errorf("destination cannot be nil")
	}

	dstVal := reflect.ValueOf(dst)

	for isPtr := dstVal.Kind() == reflect.Ptr; isPtr; isPtr = dstVal.Kind() == reflect.Ptr {
		if dstVal.IsZero() {
			dstVal.Set(reflect.New(dstVal.Type().Elem()))
		}
		dstVal = dstVal.Elem()
	}

	taggedFields, err := getTaggedFields(dstVal)
	if err != nil {
		return err
	}

	for fieldName, fieldValue := range ov.Fields {
		fieldIdx, tagged := taggedFields[fieldName]
		if !tagged {
			continue
		}

		field := dstVal.Field(fieldIdx)
		if !field.CanSet() {
			return fmt.Errorf(
				"field '%s' with tag '%s' is not accessible",
				dstVal.Type().Field(fieldIdx).Name, fieldName,
			)
		}

		for isPtr := field.Kind() == reflect.Ptr; isPtr; isPtr = field.Kind() == reflect.Ptr {
			if field.IsZero() {
				field.Set(reflect.New(field.Type().Elem()))
			}
			field = field.Elem()
		}

		var valueType reflect.Type
		if field.Kind() != reflect.Interface {
			valueType = field.Type()
		} else if valueType, err = NativeTypeForTypeHint(fieldValue.Type().Hint()); err != nil {
			return err
		}

		v := reflect.New(valueType)
		valTypeHint := TypeHintOf(v.Elem())

		if valTypeHint.TypeHint() != fieldValue.Type().Hint().TypeHint() {
			return fmt.Errorf(
				"field '%s' with tag '%s' is of uncompattible type, expected type: %s, got: %s",
				dstVal.Type().Field(fieldIdx).Name, fieldName, fieldValue.Type().Name(), v.Type().String(),
			)
		}

		if err := fieldValue.As(v.Interface()); err != nil {
			return err
		}
		field.Set(v.Elem())
	}
	return nil
}

func (ov *ObjectValue) String() string {
	return "ObjectValue"
}

func NewObjectValue(fields map[string]Valuable, fieldsTypes map[string]Typeable) Valuable {
	return &ObjectValue{Fields: fields, FieldsTypes: fieldsTypes}
}

func getTaggedFields(target reflect.Value) (map[string]int, error) {
	tagged := map[string]int{}

	var v reflect.Value = target
	for kind := target.Type().Kind(); kind == reflect.Interface || kind == reflect.Ptr; {
		elemVal := v.Elem()
		if !elemVal.IsValid() {
			break
		}
		v = elemVal
		kind = v.Type().Kind()
	}

	vType := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := vType.Field(i)
		if field.PkgPath != "" {
			continue
		}

		tag := field.Tag.Get(`cnfrmz`)
		if tag == "!" || tag == "" {
			continue
		}

		if _, alreadyTagged := tagged[tag]; alreadyTagged {
			return nil, fmt.Errorf("duplicate tag '%s'", tag)
		}
		tagged[tag] = i
	}

	return tagged, nil
}
