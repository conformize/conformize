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

type NumberIsWithinRangePredicate[T float64] struct {
	predicate.PredicateArgumentsValidator
}

func (numWithinRangePrd *NumberIsWithinRangePredicate[T]) Test(value typed.Valuable, args *typed.TupleValue) (bool, error) {
	valid, validErr := numWithinRangePrd.Validate(value, args, numWithinRangePrd.Schema())
	if !valid {
		return valid, validErr
	}

	var v T
	if err := value.As(&v); err != nil {
		return false, err
	}

	argIdx := 0
	var vLow T
	args.Elements[argIdx].As(&vLow)

	argIdx++
	var vHigh T
	args.Elements[argIdx].As(&vHigh)
	return v >= vLow && v <= vHigh, nil
}

func (numWithinRangePrd *NumberIsWithinRangePredicate[T]) Arguments() int {
	return 2
}

func (numWithinRangePrd *NumberIsWithinRangePredicate[T]) Schema() schema.Schemable {
	return &schema.Schema{
		Description: "Number is within range predicate",
		Version:     1,
		Attributes: map[string]schema.Attributeable{
			"Value": &attributes.NumberAttribute{
				Required: true,
			},
			"Arguments": &attributes.TupleAttribute{
				Required:      true,
				ElementsTypes: []typed.Typeable{&typed.NumberTyped{}, &typed.NumberTyped{}},
			},
		},
	}
}
