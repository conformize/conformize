// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package list

import (
	"time"

	"github.com/conformize/conformize/common/typed"
	"github.com/conformize/conformize/internal/providers/api/schema"
	"github.com/conformize/conformize/internal/providers/api/schema/attributes"
	"github.com/conformize/conformize/predicates/predicate"
	"github.com/spaolacci/murmur3"
)

type ListIsSubsetPredicate struct {
	predicate.PredicateArgumentsValidator
	Args *typed.TupleValue
}

func (listIsSubsetPrd *ListIsSubsetPredicate) Test(value typed.Valuable) (bool, error) {
	valid, validErr := listIsSubsetPrd.Validate(value, listIsSubsetPrd.Args, listIsSubsetPrd.Schema())
	if !valid {
		return valid, validErr
	}

	listVal := value.(*typed.ListValue)
	oListVal := listIsSubsetPrd.Args.Elements[0].(*typed.ListValue)

	if len(listVal.Elements) > len(oListVal.Elements) {
		return false, nil
	}

	seen := make(map[[16]byte]int)
	seed := uint32(time.Now().UnixNano())
	hasher := murmur3.New128WithSeed(seed)
	for _, elem := range oListVal.Elements {
		elemKey, err := elementHash(elem, hasher, seed)
		if err != nil {
			return false, err
		}
		seen[elemKey]++
	}

	for idx, elem := range listVal.Elements {
		if elem.Type().Hint() != oListVal.Elements[idx].Type().Hint() {
			return false, nil
		}

		elemKey, err := elementHash(elem, hasher, seed)
		if err != nil {
			return false, err
		}

		occurences, ok := seen[elemKey]
		if !ok {
			return false, nil
		}
		occurences--
		if 0 > occurences {
			return false, nil
		}
		seen[elemKey] = occurences
	}
	return true, nil
}

func (listIsSubsetPrd *ListIsSubsetPredicate) ArgumentsLength() int {
	return 1
}

func (listIsSubsetPrd *ListIsSubsetPredicate) Arguments(args *typed.TupleValue) {
	listIsSubsetPrd.Args = args
}

func (listIsSubsetPrd *ListIsSubsetPredicate) Schema() schema.Schemable {
	return &schema.Schema{
		Description: "List is a subset predicate",
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
