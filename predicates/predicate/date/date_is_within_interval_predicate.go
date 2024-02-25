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
)

type DateIsWithinIntervalPredicate struct {
	Args typed.Valuable
}

func (dateIsWithinIntervalPrd *DateIsWithinIntervalPredicate) Test(value typed.Valuable) (bool, error) {
	var dateVal string
	value.As(&dateVal)

	date, err := time.Parse(dateVal, dateVal)
	if err != nil {
		return false, err
	}

	argIdx := 0
	args := dateIsWithinIntervalPrd.Args.(*typed.ListValue)

	var startDateVal string
	startDateArg := args.Elements[argIdx]
	startDateArg.As(&startDateVal)

	startDate, err := time.Parse(startDateVal, startDateVal)
	if err != nil {
		return false, err
	}

	argIdx++
	var endDateVal string
	endDateArg := args.Elements[argIdx]
	endDateArg.As(&endDateVal)
	endDate, err := time.Parse(endDateVal, endDateVal)
	if err != nil {
		return false, err
	}
	return !date.Before(startDate) && !date.After(endDate), nil
}

func (dateIsWithinIntervalPrd *DateIsWithinIntervalPredicate) ArgumentsCount() int {
	return 2
}

func (dateIsWithinIntervalPrd *DateIsWithinIntervalPredicate) Arguments(args typed.Valuable) {
	dateIsWithinIntervalPrd.Args = args
}

func (dateIsWithinIntervalPrd *DateIsWithinIntervalPredicate) Schema() schema.Schemable {
	return &schema.Schema{
		Description: "Date is within interval predicate",
		Version:     1,
		Attributes: map[string]schema.Attributeable{
			"Value": &attributes.StringAttribute{
				Required: true,
			},
			"Arguments": &attributes.ListAttribute{
				Required:     true,
				ElementsType: &typed.StringTyped{},
			},
		},
	}
}
