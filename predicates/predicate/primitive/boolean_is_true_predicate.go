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
)

type BooleanIsTruePredicate struct {
	BooleanIsFalsePredicate
	Args *typed.TupleValue
}

func (boolIsTruePrd *BooleanIsTruePredicate) Test(value typed.Valuable) (bool, error) {
	res, err := boolIsTruePrd.BooleanIsFalsePredicate.Test(value)
	return err == nil && !res, err
}

func (boolIsTruePrd *BooleanIsTruePredicate) ArgumentsLength() int {
	return 0
}

func (boolIsTruePrd *BooleanIsTruePredicate) Arguments(args *typed.TupleValue) {
	boolIsTruePrd.Args = args
}

func (boolIsTruePrd *BooleanIsTruePredicate) Schema() schema.Schemable {
	return &schema.Schema{
		Description: "Boolean is true predicate",
		Version:     1,
		Attributes: map[string]schema.Attributeable{
			"Value": &attributes.BooleanAttribute{
				Required: true,
			},
		},
	}
}
