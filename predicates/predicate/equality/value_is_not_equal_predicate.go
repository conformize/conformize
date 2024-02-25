// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package equality

import (
	"github.com/conformize/conformize/common/typed"
	"github.com/conformize/conformize/internal/providers/api/schema"
	"github.com/conformize/conformize/internal/providers/api/schema/attributes"
)

type ValueIsNotEqualPredicate struct{}

func (valIsNotEqPrd *ValueIsNotEqualPredicate) Test(value typed.Valuable, args typed.Valuable) (bool, error) {
	valEqPrd := &ValueIsEqualPredicate{}
	equal, err := valEqPrd.Test(value, args)
	return !equal && err == nil, err
}

func (valIsNotEqPrd *ValueIsNotEqualPredicate) Schema() schema.Schemable {
	return &schema.Schema{
		Description: "Value equality predicate",
		Version:     1,
		Attributes: map[string]schema.Attributeable{
			"Value": &attributes.VariantAttribute{
				Required: true,
				VariantsTypes: []typed.Typeable{
					&typed.BooleanTyped{},
					&typed.NumberTyped{},
					&typed.StringTyped{},
					&typed.ListTyped{ElementsType: &typed.GenericTyped{}},
				},
			},
			"Arguments": &attributes.VariantAttribute{
				Required: true,
				VariantsTypes: []typed.Typeable{
					&typed.BooleanTyped{},
					&typed.NumberTyped{},
					&typed.StringTyped{},
					&typed.ListTyped{ElementsType: &typed.GenericTyped{}},
				},
			},
		},
	}
}
