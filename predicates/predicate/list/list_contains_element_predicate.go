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
	"github.com/conformize/conformize/predicates/predicate"
)

type ListContainsElementPredicate struct {
	predicate.PredicateArgumentsValidator
	PredicateBuilder predicates.PredicateBuilder
	Args             *typed.TupleValue
}

func (listContainsPrd *ListContainsElementPredicate) Test(value typed.Valuable) (bool, error) {
	valid, validErr := listContainsPrd.Validate(value, listContainsPrd.Args, listContainsPrd.Schema())
	if !valid {
		return valid, validErr
	}

	listVal := value.(*typed.ListValue)
	element := listContainsPrd.Args.Elements[0]
	return listContainsPrd.contains(listVal, element)
}

func (listContainsPrd *ListContainsElementPredicate) contains(listVal *typed.ListValue, element typed.Valuable) (bool, error) {
	listLen := len(listVal.Elements)
	for startIdx, endIdx := 0, listLen-1; startIdx <= endIdx; startIdx, endIdx = startIdx+1, endIdx-1 {
		for _, idx := range []int{startIdx, endIdx} {
			elemEqPrd, err := listContainsPrd.PredicateBuilder.Build(condition.EQUAL)
			if err != nil {
				return false, err
			}

			elemEqPrd.Arguments(
				&typed.TupleValue{
					Elements:      []typed.Valuable{element},
					ElementsTypes: []typed.Typeable{element.Type()},
				},
			)
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

func (listContainsPrd *ListContainsElementPredicate) ArgumentsLength() int {
	return 1
}

func (listContainsPrd *ListContainsElementPredicate) Arguments(args *typed.TupleValue) {
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
			"Arguments": &attributes.TupleAttribute{
				Required: true,
				ElementsTypes: []typed.Typeable{
					&typed.GenericTyped{},
				},
			},
		},
	}
}
