// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package collection

import (
	"fmt"

	"github.com/conformize/conformize/common/typed"
	"github.com/conformize/conformize/internal/providers/api/schema"
	"github.com/conformize/conformize/internal/providers/api/schema/attributes"
)

type IsEmptyPredicate struct{}

func (isEmptyPrd *IsEmptyPredicate) Test(value typed.Valuable, _ typed.Valuable) (bool, error) {
	elemVal, ok := value.(typed.Elementable)
	if !ok {
		return false, fmt.Errorf("invalid value type: expected list or tuple value")
	}
	return elemVal.Length() == 0, nil
}

func (isEmptyPrd *IsEmptyPredicate) Schema() schema.Schemable {
	return &schema.Schema{
		Description: "List is empty predicate",
		Version:     1,
		Attributes: map[string]schema.Attributeable{
			"Value": &attributes.ListAttribute{
				Required:     true,
				ElementsType: &typed.GenericTyped{},
			},
		},
	}
}
