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

type BooleanIsEqualPredicate struct {
	predicate.PredicateArgumentsValidator
}

func (boolEqPrd *BooleanIsEqualPredicate) Test(value typed.Valuable, args *typed.TupleValue) (bool, error) {
	valid, validErr := boolEqPrd.Validate(value, args, boolEqPrd.Schema())
	if !valid {
		return valid, validErr
	}

	var v bool
	value.As(&v)

	var vo bool
	args.Elements[0].As(&vo)
	return v == vo, nil
}

func (boolEqPrd *BooleanIsEqualPredicate) Arguments() int {
	return 1
}

func (boolEqPrd *BooleanIsEqualPredicate) Schema() schema.Schemable {
	return &schema.Schema{
		Description: "Boolean equality predicate",
		Version:     1,
		Attributes: map[string]schema.Attributeable{
			"Value": &attributes.BooleanAttribute{
				Required: true,
			},
			"Arguments": &attributes.TupleAttribute{
				Required: true,
				ElementsTypes: []typed.Typeable{
					&typed.BooleanTyped{},
				},
			},
		},
	}
}
