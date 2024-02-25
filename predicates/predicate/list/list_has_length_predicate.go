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

type ListHasLengthPredicate struct {
	predicate.PredicateArgumentsValidator
	PredicateBuilder predicates.PredicateBuilder
}

func (listLenPrd *ListHasLengthPredicate) Test(value typed.Valuable, args *typed.TupleValue) (bool, error) {
	valid, validErr := listLenPrd.Validate(value, args, listLenPrd.Schema())
	if !valid {
		return valid, validErr
	}

	listVal := value.(*typed.ListValue)
	argIdx := 0
	var conditionStr string
	args.Elements[argIdx].As(&conditionStr)

	cond := condition.FromString(conditionStr)
	supportedConditions := listLenPrd.supportedConditions()
	if _, ok := supportedConditions[cond]; !ok {
		return false, fmt.Errorf("unknown condition %s", conditionStr)
	}

	lenVal, _ := typed.NewNumberValue(len(listVal.Elements))
	condPrd, _ := listLenPrd.PredicateBuilder.Build(cond)

	argIdx++
	condArgs := &typed.TupleValue{Elements: args.Elements[argIdx:], ElementsTypes: args.ElementsTypes[argIdx:]}
	return condPrd.Test(lenVal, condArgs)
}

func (listLenPrd *ListHasLengthPredicate) Arguments() int {
	return 2
}

func (strLenPrd *ListHasLengthPredicate) supportedConditions() map[condition.ConditionType]struct{} {
	return map[condition.ConditionType]struct{}{
		condition.EQUAL:                 {},
		condition.GREATER_THAN:          {},
		condition.GREATER_THAN_OR_EQUAL: {},
		condition.LESS_THAN:             {},
		condition.LESS_THAN_OR_EQUAL:    {},
		condition.WITHIN_RANGE:          {},
	}
}

func (listLenPrd *ListHasLengthPredicate) Schema() schema.Schemable {
	return &schema.Schema{
		Description: "List length predicate",
		Version:     1,
		Attributes: map[string]schema.Attributeable{
			"Value": &attributes.ListAttribute{
				Required:     true,
				ElementsType: &typed.GenericTyped{},
			},
			"Arguments": &attributes.TupleAttribute{
				Required: true,
				ElementsTypes: []typed.Typeable{
					&typed.StringTyped{},
					&typed.NumberTyped{},
				},
			},
		},
	}
}
