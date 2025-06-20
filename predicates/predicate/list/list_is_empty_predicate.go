// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package list

import (
	"github.com/conformize/conformize/common/typed"
	"github.com/conformize/conformize/internal/providers/api/schema"
	"github.com/conformize/conformize/internal/providers/api/schema/attributes"
)

type ListIsEmptyPredicate struct {
	Args typed.Valuable
}

func (listIsEmptyPrd *ListIsEmptyPredicate) Test(value typed.Valuable) (bool, error) {
	listVal := value.(*typed.ListValue)
	return len(listVal.Elements) == 0, nil
}

func (listIsEmptyPrd *ListIsEmptyPredicate) ArgumentsCount() int {
	return 0
}

func (listIsEmptyPrd *ListIsEmptyPredicate) Arguments(args typed.Valuable) {
	listIsEmptyPrd.Args = args
}

func (listIsEmptyPrd *ListIsEmptyPredicate) Schema() schema.Schemable {
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
