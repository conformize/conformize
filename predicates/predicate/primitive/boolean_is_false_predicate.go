// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package primitive

import (
	"fmt"

	"github.com/conformize/conformize/common/typed"
	"github.com/conformize/conformize/internal/providers/api/schema"
	"github.com/conformize/conformize/internal/providers/api/schema/attributes"
)

type BooleanIsFalsePredicate struct{}

func (boolIsFalsePrd *BooleanIsFalsePredicate) Test(value typed.Valuable, _ typed.Valuable) (bool, error) {
	if value == nil {
		return false, nil
	}

	if value.Type().Hint().TypeHint() != typed.Boolean {
		return false, fmt.Errorf("expected a boolean value, got %s", value.Type().Name())
	}

	var v bool
	value.As(&v)
	return !v, nil
}

func (boolIsFalsePrd *BooleanIsFalsePredicate) Schema() schema.Schemable {
	return &schema.Schema{
		Description: "Boolean is false predicate",
		Version:     1,
		Attributes: map[string]schema.Attributeable{
			"Value": &attributes.BooleanAttribute{
				Required: true,
			},
		},
	}
}
