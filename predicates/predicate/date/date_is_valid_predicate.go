// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package date

import (
	"time"

	"github.com/conformize/conformize/common/typed"
	"github.com/conformize/conformize/internal/providers/api/schema"
	"github.com/conformize/conformize/internal/providers/api/schema/attributes"
	"github.com/conformize/conformize/predicates/predicate"
)

type DateIsValidPredicate struct {
	predicate.PredicateArgumentsValidator
}

func (dateIsValidPrd *DateIsValidPredicate) Test(value typed.Valuable, args *typed.TupleValue) (bool, error) {
	valid, validErr := dateIsValidPrd.Validate(value, args, dateIsValidPrd.Schema())
	if !valid {
		return valid, validErr
	}

	var dateVal string
	value.As(&dateVal)

	_, err := time.Parse(dateVal, dateVal)
	return err == nil, err
}

func (dateIsValidPrd *DateIsValidPredicate) Arguments() int {
	return 0
}

func (dateIsNotPrd *DateIsValidPredicate) Schema() schema.Schemable {
	return &schema.Schema{
		Description: "Date is valid predicate",
		Version:     1,
		Attributes: map[string]schema.Attributeable{
			"Value": &attributes.StringAttribute{
				Required: true,
			},
		},
	}
}