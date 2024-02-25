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

func (ov *ObjectValue) As(dst interface{}) error {
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

	for fieldName, fieldValue := range ov.Fields {
		field := dstVal.FieldByName(fieldName)

		v := reflect.New(field.Type())
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
