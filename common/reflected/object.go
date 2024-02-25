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

func Object(val interface{}, targetType typed.Typeable) (typed.Valuable, error) {
	if objectType, ok := targetType.(typed.FieldsTypeable); ok {
		if val == nil {
			return &typed.ObjectValue{
				Fields:      make(map[string]typed.Valuable),
				FieldsTypes: objectType.GetFieldsTypes(),
			}, nil
		}

		var err error
		var value typed.Valuable

		reflectVal := reflect.ValueOf(val)
		fields := make(map[string]typed.Valuable, 0)
		fieldTypes := objectType.GetFieldsTypes()

		switch reflectVal.Kind() {
		case reflect.Map:
			for _, key := range reflectVal.MapKeys() {
				keyVal := key.Interface()
				v := reflectVal.MapIndex(key).Interface()
				valType := fieldTypes[keyVal.(string)]

				if value, err := Value(v, valType); err == nil {
					fields[keyVal.(string)] = value
				}
			}
		case reflect.Struct:
			reflectValType := reflect.TypeOf(val)
			for i := 0; i < reflectVal.NumField(); i++ {
				fieldName := reflectValType.Field(i).Name
				fieldType := fieldTypes[fieldName]
				fieldVal := reflectVal.Field(i).Interface()

				if value, err = Value(fieldVal, fieldType); err == nil {
					fields[fieldName] = value
				}
			}
		}

		if err != nil {
			return nil, err
		}
		return &typed.ObjectValue{Fields: fields, FieldsTypes: fieldTypes}, nil

	}
	return nil, fmt.Errorf("cannot reflect Object as %s type - invalid type", targetType.Name())
}
