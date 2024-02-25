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

func Object(val reflect.Value, targetType typed.Typeable) (typed.Valuable, error) {
	if objectType, ok := targetType.(typed.FieldsTypeable); ok {
		if !val.IsValid() || (val.Kind() == reflect.Ptr && val.IsNil()) {
			return &typed.ObjectValue{
				Fields:      make(map[string]typed.Valuable),
				FieldsTypes: objectType.GetFieldsTypes(),
			}, nil
		}

		for val.Kind() == reflect.Ptr {
			val = val.Elem()
		}

		var err error
		var value typed.Valuable

		fields := make(map[string]typed.Valuable)
		fieldTypes := objectType.GetFieldsTypes()

		switch val.Kind() {
		case reflect.Map:
			for _, key := range val.MapKeys() {
				keyVal := key.Interface()
				v := val.MapIndex(key)
				valType, ok := fieldTypes[keyVal.(string)]
				if !ok {
					return nil, fmt.Errorf("field %s not found in object type %s", keyVal, targetType.Name())
				}

				if value, err = Value(v, valType); err == nil {
					fields[keyVal.(string)] = value
				}
			}
		case reflect.Struct:
			reflectValType := val.Type()
			for i := range val.NumField() {
				field := reflectValType.Field(i)
				tag := field.Tag.Get(`cnfrmz`)
				if tag == "!" || tag == "" {
					continue
				}

				fieldType := fieldTypes[tag]
				fieldVal := val.Field(i)
				if value, err = Value(fieldVal, fieldType); err == nil {
					fields[tag] = value
				}
			}
		default:
			return nil, fmt.Errorf("Expected value of map or struct type, got %s", val.Type())
		}

		if err != nil {
			return nil, err
		}
		return &typed.ObjectValue{Fields: fields, FieldsTypes: fieldTypes}, nil

	}
	return nil, fmt.Errorf("cannot reflect Object as %s type - invalid type", targetType.Name())
}
