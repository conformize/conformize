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

func List(val interface{}, targetType typed.Typeable) (typed.Valuable, error) {
	if elementTyped, ok := targetType.(typed.ElementTypeable); ok {
		if val == nil {
			return &typed.ListValue{ElementsType: elementTyped.ElementType(), Elements: []typed.Valuable{}}, nil
		}

		reflectVal := reflect.ValueOf(val)
		elements := make([]typed.Valuable, 0)
		for i := 0; i < reflectVal.Len(); i++ {
			elementVal := reflectVal.Index(i).Interface()
			elemTypeHint := typed.TypeHintOf(elementVal)

			v, err := ValueFromTypeHint(elementVal, elemTypeHint)
			if err != nil {
				return nil, err
			}
			elements = append(elements, v)
		}
		return &typed.ListValue{ElementsType: elementTyped.ElementType(), Elements: elements}, nil

	}
	return nil, fmt.Errorf("invalid type: %s", targetType.Name())
}