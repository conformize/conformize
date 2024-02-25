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

type ListNoElementMatcheConditionPredicate struct {
	predicate.PredicateArgumentsValidator
	PredicateBuilder   predicates.PredicateBuilder
	conditionPredicate predicates.Predicate
}

func (listNoElemMatchesConditionPrd *ListNoElementMatcheConditionPredicate) Test(value typed.Valuable, args *typed.TupleValue) (bool, error) {
	valid, validErr := listNoElemMatchesConditionPrd.Validate(value, args, listNoElemMatchesConditionPrd.Schema())
	if !valid {
		return valid, validErr
	}

	listVal := value.(*typed.ListValue)

	argIdx := 0
	conditionVal := args.Elements[argIdx].(*typed.StringValue)

	var conditionStr string
	conditionVal.As(&conditionStr)

	cond := condition.FromString(conditionStr)
	if cond == condition.UNKNOWN {
		return false, fmt.Errorf("unknown condition %s", conditionStr)
	}

	listValLen := len(listVal.Elements)
	if listValLen == 0 {
		return false, fmt.Errorf("no elements to match. list is empty")
	}

	argIdx++
	condArgs := &typed.TupleValue{Elements: args.Elements[argIdx:], ElementsTypes: args.ElementsTypes[argIdx:]}
	return listNoElemMatchesConditionPrd.noMatch(listVal, cond, condArgs)
}

func (listNoElemMatchesConditionPrd *ListNoElementMatcheConditionPredicate) noMatch(listVal *typed.ListValue, condition condition.ConditionType, condArgs *typed.TupleValue) (bool, error) {
	for startIdx, endIdx := 0, len(listVal.Elements)-1; endIdx >= startIdx; startIdx, endIdx = startIdx+1, endIdx-1 {
		for _, idx := range []int{startIdx, endIdx} {
			element := listVal.Elements[idx]
			if listNoElemMatchesConditionPrd.conditionPredicate == nil {
				prd, err := listNoElemMatchesConditionPrd.PredicateBuilder.Build(condition)
				if err != nil {
					return false, err
				}
				listNoElemMatchesConditionPrd.conditionPredicate = prd
			}

			ok, err := listNoElemMatchesConditionPrd.conditionPredicate.Test(element, condArgs)
			if ok || err != nil {
				return false, err
			}
			if startIdx == endIdx {
				break
			}
		}
	}
	return true, nil
}

func (listNoElemMatchesConditionPrd *ListNoElementMatcheConditionPredicate) Arguments() int {
	return 1
}

func (listNoElemMatchesConditionPrd *ListNoElementMatcheConditionPredicate) Schema() schema.Schemable {
	return &schema.Schema{
		Description: "No element in a list meets condition predicate",
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
				},
			},
		},
	}
}
