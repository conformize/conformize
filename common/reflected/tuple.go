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

func Tuple(val interface{}, targetType typed.Typeable) (typed.Valuable, error) {
	if tupleTyped, ok := targetType.(typed.MixedElementsTypeable); ok {
		if val == nil {
			return nil, fmt.Errorf("cannot reflect nil as %s type", tupleTyped.Name())
		}

		reflectVal := reflect.ValueOf(val)
		if reflectVal.Kind() == reflect.Slice {
			sliceLen := reflectVal.Len()
			elemTypes := tupleTyped.GetElementsTypes()
			if sliceLen != len(elemTypes) {
				return nil, fmt.Errorf("cannot reflect as %s type - value's elements count doesn't match", targetType.Name())
			}

			tupleElements := make([]typed.Valuable, sliceLen)
			for idx := 0; idx < sliceLen; idx++ {
				reflectElemVal := reflectVal.Index(idx).Interface()
				elemTypeHint := typed.TypeHintOf(reflectElemVal)
				val, err := ValueFromTypeHint(reflectElemVal, elemTypeHint)
				if err != nil {
					return nil, err
				}
				tupleElements[idx] = val
			}
			return &typed.TupleValue{Elements: tupleElements, ElementsTypes: elemTypes}, nil
		}
	}
	return nil, fmt.Errorf("cannot reflect as %s type - invalid value type", targetType.Name())
}