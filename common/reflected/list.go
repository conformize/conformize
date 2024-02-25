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

func List(val reflect.Value, targetType typed.Typeable) (typed.Valuable, error) {
	if elementTyped, ok := targetType.(typed.ElementTypeable); ok {
		if !val.IsValid() || (val.Kind() == reflect.Ptr && val.IsNil()) {
			return &typed.ListValue{ElementsType: elementTyped.ElementType(), Elements: []typed.Valuable{}}, nil
		}

		elements := make([]typed.Valuable, val.Len())
		elemType := elementTyped.ElementType()

		var elementVal reflect.Value
		var v typed.Valuable
		var err error
		for i := range val.Len() {
			elementVal = val.Index(i)
			v, err = Value(elementVal, elemType)
			if err != nil {
				return nil, err
			}
			elements[i] = v
		}
		return &typed.ListValue{ElementsType: elementTyped.ElementType(), Elements: elements}, nil

	}
	return nil, fmt.Errorf("invalid type: %s", targetType.Name())
}
