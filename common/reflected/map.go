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

func Map(val reflect.Value, targetType typed.Typeable) (typed.Valuable, error) {
	if elementTyped, ok := targetType.(typed.ElementTypeable); ok {
		if !val.IsValid() || (val.Kind() == reflect.Ptr && val.IsNil()) {
			return &typed.MapValue{
				ElementsType: elementTyped.ElementType(),
				Elements:     make(map[string]typed.Valuable),
			}, nil
		}

		elemType := elementTyped.ElementType()
		elements := make(map[string]typed.Valuable)

		iter := val.MapRange()
		var v typed.Valuable
		var err error
		for iter.Next() {
			key := iter.Key().String()
			v, err = Value(iter.Value(), elemType)
			if err != nil {
				return nil, err
			}
			elements[key] = v
		}
		return &typed.MapValue{ElementsType: elementTyped.ElementType(), Elements: elements}, nil
	}
	return nil, fmt.Errorf("invalid type: %s", targetType.Name())
}
