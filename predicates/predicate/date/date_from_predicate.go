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

type DateFromPredicate struct {
	predicate.PredicateArgumentsValidator
}

func (dateFromPredicate *DateFromPredicate) Test(value typed.Valuable, args *typed.TupleValue) (bool, error) {
	valid, validErr := dateFromPredicate.Validate(value, args, dateFromPredicate.Schema())
	if !valid {
		return valid, validErr
	}

	var dateVal string
	value.As(&dateVal)

	date, err := time.Parse(dateVal, dateVal)
	if err != nil {
		return false, err
	}

	var oDateVal string
	oDateArg := args.Elements[0]
	oDateArg.As(&oDateVal)

	oDate, err := time.Parse(oDateVal, oDateVal)
	if err != nil {
		return false, err
	}
	return date.After(oDate) || date.Equal(oDate), nil
}

func (dateFromPredicate *DateFromPredicate) Arguments() int {
	return 1
}

func (dateFromPredicate *DateFromPredicate) Schema() schema.Schemable {
	return &schema.Schema{
		Description: "Date is equal or after predicate",
		Version:     1,
		Attributes: map[string]schema.Attributeable{
			"Value": &attributes.StringAttribute{
				Required: true,
			},
			"Arguments": &attributes.TupleAttribute{
				Required:      true,
				ElementsTypes: []typed.Typeable{&typed.StringTyped{}},
			},
		},
	}
}