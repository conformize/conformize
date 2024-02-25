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
	"github.com/conformize/conformize/predicates"
	"github.com/conformize/conformize/predicates/predicate"
	"github.com/spaolacci/murmur3"
)

type ListContainsAllElementsPredicate struct {
	predicate.PredicateArgumentsValidator
	PredicateBuilder predicates.PredicateBuilder
}

func (listContainsAllPrd *ListContainsAllElementsPredicate) Test(value typed.Valuable, args *typed.TupleValue) (bool, error) {
	valid, validErr := listContainsAllPrd.Validate(value, args, listContainsAllPrd.Schema())
	if !valid {
		return valid, validErr
	}

	listVal := value.(*typed.ListValue)
	elementsListVal := args.Elements[0].(*typed.ListValue)
	matches := make(map[[16]byte]bool)
	seed := uint32(time.Now().UnixNano())
	hasher := murmur3.New128WithSeed(seed)

	for _, elem := range elementsListVal.Elements {
		elemKey, err := elementHash(elem, hasher, seed)
		if err != nil {
			return false, err
		}
		matches[elemKey] = false
	}

	unmatchedCount := len(matches)
	for startIdx, endIdx := 0, len(listVal.Elements)-1; startIdx <= endIdx; startIdx, endIdx = startIdx+1, endIdx-1 {
		for _, idx := range []int{startIdx, endIdx} {
			elemHash, err := elementHash(listVal.Elements[idx], hasher, seed)
			if err != nil {
				return false, err
			}

			if matched, ok := matches[elemHash]; ok && !matched {
				matches[elemHash] = true
				unmatchedCount--
			}

			if unmatchedCount == 0 || startIdx == endIdx {
				break
			}
		}
	}
	return unmatchedCount == 0, nil
}

func (listContainsAllPrd *ListContainsAllElementsPredicate) Arguments() int {
	return 1
}

func (listContainsAllPrd *ListContainsAllElementsPredicate) Schema() schema.Schemable {
	return &schema.Schema{
		Description: "All elements are in a list predicate",
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
