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

type NumberIsWithinRangePredicate[T float64] struct{}

func (numWithinRangePrd *NumberIsWithinRangePredicate[T]) Test(value typed.Valuable, args typed.Valuable) (bool, error) {
	if value == nil {
		return false, fmt.Errorf("value is nil")
	}

	if value.Type().Hint().TypeHint() != typed.Number {
		return false, fmt.Errorf("expected a number value, got %s", value.Type().Name())
	}

	var v T
	if err := value.As(&v); err != nil {
		return false, err
	}

	if args == nil {
		return false, fmt.Errorf("arguments are nil")
	}

	if args.Type().Hint().TypeHint() != typed.List {
		return false, fmt.Errorf("expected a list of numbers as arguments, got %s", args.Type().Name())
	}

	argIdx := 0
	listArg, ok := args.(*typed.ListValue)
	if !ok {
		return false, fmt.Errorf("expected a list of numbers as arguments, got %s", args.Type().Name())
	}

	if len(listArg.Elements) != 2 {
		return false, fmt.Errorf("expected exactly 2 arguments for range, got %d", len(listArg.Elements))
	}

	if listArg.Elements[0].Type().Hint().TypeHint() != typed.Number {
		return false, fmt.Errorf("expected a number as the first argument, got %s", listArg.Elements[0].Type().Name())
	}

	var vLow T
	listArg.Elements[argIdx].As(&vLow)

	argIdx++

	if listArg.Elements[argIdx].Type().Hint().TypeHint() != typed.Number {
		return false, fmt.Errorf("expected a number as the second argument, got %s", listArg.Elements[argIdx].Type().Name())
	}

	var vHigh T
	listArg.Elements[argIdx].As(&vHigh)
	return v >= vLow && v <= vHigh, nil
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
