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
	"github.com/conformize/conformize/predicates/predicate"
)

type ListIsEqualPredicate struct {
	predicate.PredicateArgumentsValidator
	PredicateBuilder predicates.PredicateBuilder
}

func (listEqPrd *ListIsEqualPredicate) Test(value typed.Valuable, args *typed.TupleValue) (bool, error) {
	valid, validErr := listEqPrd.Validate(value, args, listEqPrd.Schema())
	if !valid {
		return valid, validErr
	}

	listVal := value.(*typed.ListValue)
	oListVal := args.Elements[0].(*typed.ListValue)

	if len(listVal.Elements) != len(oListVal.Elements) {
		return false, nil
	}

	if listVal.ElementsType.Hint() != oListVal.ElementsType.Hint() {
		return false, fmt.Errorf("cannot compare lists with elements of different type")
	}

	listLen := len(oListVal.Elements)
	if listLen == 0 {
		return true, nil
	}

	for startIdx, endIdx := 0, len(listVal.Elements)-1; endIdx >= startIdx; startIdx, endIdx = startIdx+1, endIdx-1 {
		for _, idx := range []int{startIdx, endIdx} {
			elemEqPrd, err := listEqPrd.PredicateBuilder.Build(condition.EQUAL)
			if err != nil {
				return false, err
			}

			ok, err := elemEqPrd.Test(listVal.Elements[idx], &typed.TupleValue{
				Elements:      []typed.Valuable{oListVal.Elements[idx]},
				ElementsTypes: []typed.Typeable{oListVal.Elements[idx].Type()},
			})

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

func (listEqPrd *ListIsEqualPredicate) Arguments() int {
	return 1
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
