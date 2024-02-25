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

type BooleanIsFalsePredicate struct {
	predicate.PredicateArgumentsValidator
}

func (boolIsFalsePrd *BooleanIsFalsePredicate) Test(value typed.Valuable, args *typed.TupleValue) (bool, error) {
	valid, validErr := boolIsFalsePrd.Validate(value, args, boolIsFalsePrd.Schema())
	if !valid {
		return valid, validErr
	}

	var v bool
	value.As(&v)
	return !v, nil
}

func (boolIsFalsePrd *BooleanIsFalsePredicate) Arguments() int {
	return 0
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
