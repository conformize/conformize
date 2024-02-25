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

func Map(val interface{}, targetType typed.Typeable) (typed.Valuable, error) {
	if elementTyped, ok := targetType.(typed.ElementTypeable); ok {
		if val == nil {
			return &typed.MapValue{
				ElementsType: elementTyped.ElementType(),
				Elements:     make(map[string]typed.Valuable),
			}, nil
		}

		reflectVal := reflect.ValueOf(val)
		elemType := elementTyped.ElementType()
		elements := make(map[string]typed.Valuable)

		iter := reflectVal.MapRange()
		for iter.Next() {
			key := iter.Key().String()
			v := iter.Value().Interface()
			val, err := Value(v, elemType)
			if err != nil {
				return nil, err
			}
			elements[key] = val
		}
		return &typed.MapValue{ElementsType: elementTyped.ElementType(), Elements: elements}, nil
	}
	return nil, fmt.Errorf("invalid type: %s", targetType.Name())
}
