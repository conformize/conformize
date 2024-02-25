// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package list

import (
	"fmt"

	"github.com/conformize/conformize/common/typed"
	"github.com/conformize/conformize/internal/providers/api/schema"
	"github.com/conformize/conformize/internal/providers/api/schema/attributes"
	"github.com/conformize/conformize/predicates"
	"github.com/conformize/conformize/predicates/condition"
)

type ListIsEqualPredicate struct {
	PredicateBuilder predicates.PredicateBuilder
	Args             typed.Valuable
}

func (listEqPrd *ListIsEqualPredicate) Test(value typed.Valuable) (bool, error) {
	if listEqPrd.Args == nil || value == nil {
		return false, fmt.Errorf("arguments or value cannot be nil")
	}

	listVal, ok := value.(*typed.ListValue)
	if !ok {
		return false, fmt.Errorf("expected a list value, got %s", value.Type().Name())
	}

	oListVal, ok := listEqPrd.Args.(*typed.ListValue)
	if !ok {
		return false, fmt.Errorf("expected a list value as argument, got %s", listEqPrd.Args.Type().Name())
	}

	if len(listVal.Elements) != len(oListVal.Elements) {
		return false, nil
	}

	if len(listVal.Elements) == 0 && len(oListVal.Elements) == 0 {
		return true, nil
	}

	if listVal.ElementsType.Hint() != oListVal.ElementsType.Hint() {
		return false, fmt.Errorf("cannot compare lists with elements of different type")
	}

	for startIdx, endIdx := 0, len(listVal.Elements)-1; endIdx >= startIdx; startIdx, endIdx = startIdx+1, endIdx-1 {
		for _, idx := range []int{startIdx, endIdx} {
			prd, err := listEqPrd.PredicateBuilder.Build(condition.EQ)
			if err != nil {
				return false, err
			}

			elemEqPrd := prd.(predicates.ArgumentsPredicate)
			elemEqPrd.Arguments(oListVal.Elements[idx])
			ok, err := elemEqPrd.Test(listVal.Elements[idx])

			if !ok || err != nil {
				return false, err
			}
			if startIdx == endIdx {
				break
			}
		}
	}
	return true, nil
}

func (listEqPrd *ListIsEqualPredicate) ArgumentsCount() int {
	return 1
}

func (listEqPrd *ListIsEqualPredicate) Arguments(args typed.Valuable) {
	listEqPrd.Args = args
}

func (listEqPrd *ListIsEqualPredicate) Schema() schema.Schemable {
	return &schema.Schema{
		Description: "List equality predicate",
		Version:     1,
		Attributes: map[string]schema.Attributeable{
			"Value": &attributes.ListAttribute{
				Required:     true,
				ElementsType: &typed.GenericTyped{},
			},
			"Arguments": &attributes.TupleAttribute{
				Required: true,
				ElementsTypes: []typed.Typeable{
					&typed.ListTyped{ElementsType: &typed.GenericTyped{}},
				},
			},
		},
	}
}
