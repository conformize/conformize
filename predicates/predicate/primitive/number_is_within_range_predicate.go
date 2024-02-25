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

type NumberIsWithinRangePredicate[T float64] struct {
	Args typed.Valuable
}

func (numWithinRangePrd *NumberIsWithinRangePredicate[T]) Test(value typed.Valuable) (bool, error) {
	if value == nil {
		return false, fmt.Errorf("value is nil")
	}

	if value.Type().Hint() != typed.Number {
		return false, fmt.Errorf("expected a number value, got %s", value.Type().Name())
	}

	var v T
	if err := value.As(&v); err != nil {
		return false, err
	}

	if numWithinRangePrd.Args == nil {
		return false, fmt.Errorf("arguments are nil")
	}

	if numWithinRangePrd.Args.Type().Hint() != typed.List {
		return false, fmt.Errorf("expected a list of numbers as arguments, got %s", numWithinRangePrd.Args.Type().Name())
	}

	argIdx := 0
	args := numWithinRangePrd.Args.(*typed.ListValue)

	if len(args.Elements) != 2 {
		return false, fmt.Errorf("expected exactly 2 arguments for range, got %d", len(args.Elements))
	}

	if args.Elements[0].Type().Hint() != typed.Number {
		return false, fmt.Errorf("expected a number as the first argument, got %s", args.Elements[0].Type().Name())
	}

	var vLow T
	args.Elements[argIdx].As(&vLow)

	argIdx++

	if args.Elements[argIdx].Type().Hint() != typed.Number {
		return false, fmt.Errorf("expected a number as the second argument, got %s", args.Elements[argIdx].Type().Name())
	}

	var vHigh T
	args.Elements[argIdx].As(&vHigh)
	return v >= vLow && v <= vHigh, nil
}

func (numWithinRangePrd *NumberIsWithinRangePredicate[T]) ArgumentsCount() int {
	return 2
}

func (numWithinRangePrd *NumberIsWithinRangePredicate[T]) Arguments(args typed.Valuable) {
	numWithinRangePrd.Args = args
}

func (numWithinRangePrd *NumberIsWithinRangePredicate[T]) Schema() schema.Schemable {
	return &schema.Schema{
		Description: "Number is within range predicate",
		Version:     1,
		Attributes: map[string]schema.Attributeable{
			"Value": &attributes.NumberAttribute{
				Required: true,
			},
			"Arguments": &attributes.ListAttribute{
				Required:     true,
				ElementsType: &typed.NumberTyped{},
			},
		},
	}
}
