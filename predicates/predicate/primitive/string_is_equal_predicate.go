// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package primitive

import (
	"github.com/conformize/conformize/common/typed"
	"github.com/conformize/conformize/internal/providers/api/schema"
	"github.com/conformize/conformize/internal/providers/api/schema/attributes"
	"github.com/conformize/conformize/predicates/predicate"
)

type StringIsEqualPredicate struct {
	predicate.PredicateArgumentsValidator
}

func (strIsEqPrd *StringIsEqualPredicate) Test(value typed.Valuable, args *typed.TupleValue) (bool, error) {
	valid, validErr := strIsEqPrd.Validate(value, args, strIsEqPrd.Schema())
	if !valid {
		return valid, validErr
	}

	var s string
	value.As(&s)

	var so string
	args.Elements[0].As(&so)

	return s == so, nil
}

func (strLenPrd *StringIsEqualPredicate) Arguments() int {
	return 1
}

func (strLenPrd *StringIsEqualPredicate) Schema() schema.Schemable {
	return &schema.Schema{
		Description: "String equality predicate",
		Version:     1,
		Attributes: map[string]schema.Attributeable{
			"Value": &attributes.StringAttribute{
				Required: true,
			},
			"Arguments": &attributes.TupleAttribute{
				Required: true,
				ElementsTypes: []typed.Typeable{
					&typed.StringTyped{},
				},
			},
		},
	}
}
