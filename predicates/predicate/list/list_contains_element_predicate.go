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
	"github.com/conformize/conformize/predicates"
	"github.com/conformize/conformize/predicates/condition"
)

type ListContainsElementPredicate struct {
	PredicateBuilder predicates.PredicateBuilder
	Args             typed.Valuable
}

func (listContainsPrd *ListContainsElementPredicate) Test(value typed.Valuable) (bool, error) {
	listVal := value.(*typed.ListValue)
	element := listContainsPrd.Args
	return listContainsPrd.contains(listVal, element)
}

func (listContainsPrd *ListContainsElementPredicate) contains(listVal *typed.ListValue, element typed.Valuable) (bool, error) {
	listLen := len(listVal.Elements)
	for startIdx, endIdx := 0, listLen-1; startIdx <= endIdx; startIdx, endIdx = startIdx+1, endIdx-1 {
		for _, idx := range []int{startIdx, endIdx} {
			prd, err := listContainsPrd.PredicateBuilder.Build(condition.EQ)
			if err != nil {
				return false, err
			}

			elemEqPrd := prd.(predicates.ArgumentsPredicate)
			elemEqPrd.Arguments(element)
			ok, err := elemEqPrd.Test(listVal.Elements[idx])
			if err != nil {
				return false, err
			}
			if ok {
				return true, nil
			}
			if startIdx == endIdx {
				break
			}
		}
	}
	return false, nil
}

func (listContainsPrd *ListContainsElementPredicate) ArgumentsCount() int {
	return 1
}

func (listContainsPrd *ListContainsElementPredicate) Arguments(args typed.Valuable) {
	listContainsPrd.Args = args
}

func (listContainsPrd *ListContainsElementPredicate) Schema() schema.Schemable {
	return &schema.Schema{
		Description: "Element is in a list predicate",
		Version:     1,
		Attributes: map[string]schema.Attributeable{
			"Value": &attributes.ListAttribute{
				Required:     true,
				ElementsType: &typed.GenericTyped{},
			},
			"Arguments": &attributes.ListAttribute{
				Required:     true,
				ElementsType: &typed.GenericTyped{},
			},
		},
	}
}
