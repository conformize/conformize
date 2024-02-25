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

type NumberIsLessThanPredicate[T float64] struct {
	predicate.PredicateArgumentsValidator
}

func (numLtPrd *NumberIsLessThanPredicate[T]) Test(value typed.Valuable, args *typed.TupleValue) (bool, error) {
	valid, validErr := numLtPrd.Validate(value, args, numLtPrd.Schema())
	if !valid {
		return valid, validErr
	}

	var v T
	value.As(&v)

	var vo T
	args.Elements[0].As(&vo)
	return vo > v, nil
}

func (numLtPrd *NumberIsLessThanPredicate[T]) Arguments() int {
	return 1
}

func (numLtPrd *NumberIsLessThanPredicate[T]) Schema() schema.Schemable {
	return &schema.Schema{
		Description: "Number is less than predicate",
		Version:     1,
		Attributes: map[string]schema.Attributeable{
			"Value": &attributes.NumberAttribute{
				Required: true,
			},
			"Arguments": &attributes.TupleAttribute{
				Required: true,
				ElementsTypes: []typed.Typeable{
					&typed.NumberTyped{},
				},
			},
		},
	}
}
