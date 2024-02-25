// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package equality

import (
	"fmt"

	"github.com/conformize/conformize/common/typed"
	"github.com/conformize/conformize/internal/providers/api/schema"
	"github.com/conformize/conformize/internal/providers/api/schema/attributes"
	"github.com/conformize/conformize/predicates"
	"github.com/conformize/conformize/predicates/predicate/list"
	"github.com/conformize/conformize/predicates/predicate/primitive"
)

type ValueIsEqualPredicate struct {
	PredicateBuilder predicates.PredicateBuilder
}

func (valIsEqPrd *ValueIsEqualPredicate) Test(value typed.Valuable, args *typed.TupleValue) (bool, error) {
	var prd predicates.Predicate
	switch value.Type().Hint() {
	case typed.Boolean:
		prd = &primitive.BooleanIsEqualPredicate{}
	case typed.Number:
		prd = &primitive.NumberIsEqualPredicate[float64]{}
	case typed.String:
		prd = &primitive.StringIsEqualPredicate{}
	case typed.List:
		prd = &list.ListIsEqualPredicate{PredicateBuilder: valIsEqPrd.PredicateBuilder}
	default:
		return false, fmt.Errorf("value of type %s is not supported", value.Type().Name())
	}
	return prd.Test(value, args)
}

func (valIsEqPrd *ValueIsEqualPredicate) Arguments() int {
	return 1
}

func (valIsEqPrd *ValueIsEqualPredicate) Schema() schema.Schemable {
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
