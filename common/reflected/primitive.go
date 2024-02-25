// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package reflected

import (
	"fmt"

	"github.com/conformize/conformize/common/typed"
)

type primitiveFn func(interface{}) (typed.Valuable, error)

var primitives = []primitiveFn{
	typed.NewBooleanValue,
	typed.NewNumberValue,
	typed.NewStringValue,
}

func Primitive(value interface{}, targetType typed.Typeable) (typed.Valuable, error) {
	typeHint := targetType.Hint()
	if value != nil {
		valueTypeHint := typed.TypeHintOf(value)
		if !typed.IsPrimitive(typeHint) || valueTypeHint != typeHint {
			return nil, fmt.Errorf("invalid primitive type: %s", targetType.Name())
		}
	}
	primitiveFn := primitives[typeHint-1]
	return primitiveFn(value)
}
