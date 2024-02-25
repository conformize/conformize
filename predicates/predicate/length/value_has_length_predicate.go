// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package length

import (
	"fmt"

	"github.com/conformize/conformize/common/typed"
	"github.com/conformize/conformize/internal/providers/api/schema"
	"github.com/conformize/conformize/internal/providers/api/schema/attributes"
	"github.com/conformize/conformize/predicates"
	"github.com/conformize/conformize/predicates/predicate/list"
	"github.com/conformize/conformize/predicates/predicate/primitive"
)

type ValueHasLengthPredicate struct{}

func (strLenPrd *ValueHasLengthPredicate) Test(value typed.Valuable, args *typed.TupleValue) (bool, error) {
	var prd predicates.Predicate
	switch value.Type().Hint() {
	case typed.String:
		prd = &primitive.StringHasLengthPredicate{}
	case typed.List:
		prd = &list.ListHasLengthPredicate{}
	default:
		return false, fmt.Errorf("value of type %s is not supported", value.Type().Name())
	}
	return prd.Test(value, args)
}

func (strLenPrd *ValueHasLengthPredicate) Arguments() int {
	return 2
}

func (strLenPrd *ValueHasLengthPredicate) Schema() schema.Schemable {
	return &schema.Schema{
		Description: "Value has length predicate",
		Version:     1,
		Attributes: map[string]schema.Attributeable{
			"Value": &attributes.StringAttribute{
				Required: true,
			},
			"Arguments": &attributes.TupleAttribute{
				Required: true,
				ElementsTypes: []typed.Typeable{
					&typed.StringTyped{},
					&typed.NumberTyped{},
				},
			},
		},
	}
}
