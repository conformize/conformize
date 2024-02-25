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

type BooleanIsEqualPredicate struct {
	Args typed.Valuable
}

func (boolEqPrd *BooleanIsEqualPredicate) Test(value typed.Valuable) (bool, error) {
	if value == nil || boolEqPrd.Args == nil {
		return false, fmt.Errorf("value and arguments cannot be nil")
	}

	if value.Type().Hint() != typed.Boolean {
		return false, fmt.Errorf("expected a boolean value, got %s", value.Type().Name())
	}

	if boolEqPrd.Args.Type().Hint() != typed.Boolean {
		return false, fmt.Errorf("expected a boolean argument, got %s", boolEqPrd.Args.Type().Name())
	}

	var v bool
	value.As(&v)

	var vo bool
	boolEqPrd.Args.As(&vo)
	return v == vo, nil
}

func (boolEqPrd *BooleanIsEqualPredicate) ArgumentsCount() int {
	return 1
}

func (boolEqPrd *BooleanIsEqualPredicate) Arguments(args typed.Valuable) {
	boolEqPrd.Args = args
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
