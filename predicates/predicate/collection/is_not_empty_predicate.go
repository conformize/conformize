// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package collection

import (
	"github.com/conformize/conformize/common/typed"
	"github.com/conformize/conformize/internal/providers/api/schema"
	"github.com/conformize/conformize/internal/providers/api/schema/attributes"
)

type IsNotEmptyPredicate struct{}

func (isNotEmptyPrd *IsNotEmptyPredicate) Test(value typed.Valuable, _ typed.Valuable) (bool, error) {
	listEmptyprd := IsEmptyPredicate{}
	empty, err := listEmptyprd.Test(value, nil)
	return !empty && err == nil, err
}

func (isNotEmptyPrd *IsNotEmptyPredicate) Schema() schema.Schemable {
	return &schema.Schema{
		Description: "List is not empty predicate",
		Version:     1,
		Attributes: map[string]schema.Attributeable{
			"Value": &attributes.ListAttribute{
				Required:     true,
				ElementsType: &typed.GenericTyped{},
			},
		},
	}
}
