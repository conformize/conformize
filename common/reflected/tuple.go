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

func Tuple(val reflect.Value, targetType typed.Typeable) (typed.Valuable, error) {
	if tupleTyped, ok := targetType.(typed.MixedElementsTypeable); ok {
		if !val.IsValid() || (val.Kind() == reflect.Ptr && val.IsNil()) {
			return nil, fmt.Errorf("cannot reflect nil as %s type", tupleTyped.Name())
		}

		if val.Kind() == reflect.Slice {
			sliceLen := val.Len()
			elemTypes := tupleTyped.GetElementsTypes()
			if sliceLen != len(elemTypes) {
				return nil, fmt.Errorf("cannot reflect as %s type - value's elements count doesn't match", targetType.Name())
			}

			tupleElements := make([]typed.Valuable, sliceLen)

			var err error
			var v typed.Valuable
			var reflectElemVal reflect.Value
			for idx := range sliceLen {
				reflectElemVal = val.Index(idx)
				elemTypeHint := typed.TypeHintOf(reflectElemVal)
				v, err = ValueFromTypeHint(reflectElemVal, elemTypeHint)
				if err != nil {
					return nil, err
				}
				tupleElements[idx] = v
			}
			return &typed.TupleValue{Elements: tupleElements, ElementsTypes: elemTypes}, nil
		}
	}
	return nil, fmt.Errorf("cannot reflect as %s type - invalid value type", targetType.Name())
}
