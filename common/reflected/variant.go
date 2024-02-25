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

func Variant(val interface{}, targetType typed.Typeable) (typed.Valuable, error) {
	if variantTyped, ok := targetType.(typed.VariantTypeable); ok {
		variantTypes := variantTyped.GetVariantsTypes()
		for _, variantType := range variantTypes {
			if val, err := Value(val, variantType); err == nil {
				return &typed.VariantValue{Value: val, VariantsTypes: variantTypes}, nil
			}
		}
	}
	return nil, fmt.Errorf("cannot reflect as %s type - invalid value type", targetType.Name())
}
