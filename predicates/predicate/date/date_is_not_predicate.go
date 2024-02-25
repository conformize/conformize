// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package date

import (
	"github.com/conformize/conformize/common/typed"
	"github.com/conformize/conformize/internal/providers/api/schema"
	"github.com/conformize/conformize/internal/providers/api/schema/attributes"
)

type DateIsNotPredicate struct{}

func (dateIsNotPrd *DateIsNotPredicate) Test(value typed.Valuable, args typed.Valuable) (bool, error) {
	dateIsPrd := &DateIsPredicate{}
	same, err := dateIsPrd.Test(value, args)
	return !same && err == nil, err
}

func (dateIsNotPrd *DateIsNotPredicate) Schema() schema.Schemable {
	return &schema.Schema{
		Description: "Date is not equal predicate",
		Version:     1,
		Attributes: map[string]schema.Attributeable{
			"Value": &attributes.StringAttribute{
				Required: true,
			},
			"Arguments": &attributes.StringAttribute{
				Required: true,
			},
		},
	}
}
