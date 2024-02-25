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

type NumberIsGreaterThanOrEqualPredicate[T float64] struct {
	Args typed.Valuable
}

func (numGtOrEqPrd *NumberIsGreaterThanOrEqualPredicate[T]) Test(value typed.Valuable) (bool, error) {
	if value == nil {
		return false, fmt.Errorf("value is nil")
	}

	if numGtOrEqPrd.Args == nil {
		return false, fmt.Errorf("argument is nil")
	}

	if value.Type().Hint() != typed.Number {
		return false, fmt.Errorf("expected a number value, got %s", value.Type().Name())
	}

	if numGtOrEqPrd.Args.Type().Hint() != typed.Number {
		return false, fmt.Errorf("expected a number argument, got %s", numGtOrEqPrd.Args.Type().Name())
	}

	var v T
	value.As(&v)

	var vo T
	numGtOrEqPrd.Args.As(&vo)
	return v >= vo, nil
}

func (numGtOrEqPrd *NumberIsGreaterThanOrEqualPredicate[T]) ArgumentsCount() int {
	return 1
}

func (numGtOrEqPrd *NumberIsGreaterThanOrEqualPredicate[T]) Arguments(args typed.Valuable) {
	numGtOrEqPrd.Args = args
}

func (numGtOrEqPrd *NumberIsGreaterThanOrEqualPredicate[T]) Schema() schema.Schemable {
	return &schema.Schema{
		Description: "Number is greater than or equal predicate",
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
