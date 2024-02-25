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
	"github.com/spaolacci/murmur3"
)

type ListContainsAnyOfElementsPredicate struct {
	PredicateBuilder predicates.PredicateBuilder
	Args             typed.Valuable
}

func (listContainsAnyOfElemPrd *ListContainsAnyOfElementsPredicate) Test(value typed.Valuable) (bool, error) {
	listVal := value.(*typed.ListValue)
	elementsListVal := listContainsAnyOfElemPrd.Args.(*typed.ListValue)

	seen := make(map[[16]byte]struct{})
	seed := uint32(time.Now().UnixNano())
	hasher := murmur3.New128WithSeed(seed)
	for _, elem := range elementsListVal.Elements {
		elemKey, err := elementHash(elem, hasher, seed)
		if err != nil {
			return false, err
		}
		seen[elemKey] = struct{}{}
	}

	for startIdx, endIdx := 0, len(listVal.Elements)-1; startIdx <= endIdx; startIdx, endIdx = startIdx+1, endIdx-1 {
		for _, idx := range []int{startIdx, endIdx} {
			elemHash, err := elementHash(listVal.Elements[idx], hasher, seed)
			if err != nil {
				return false, err
			}

			if _, ok := seen[elemHash]; ok {
				return true, nil
			}

			if startIdx == endIdx {
				break
			}
		}
	}
	return false, nil
}

func (listContainsAnyPrd *ListContainsAnyOfElementsPredicate) ArgumentsCount() int {
	return 1
}

func (listContainsAnyPrd *ListContainsAnyOfElementsPredicate) Arguments(args typed.Valuable) {
	listContainsAnyPrd.Args = args
}

func (listContainsAnyPrd *ListContainsAnyOfElementsPredicate) Schema() schema.Schemable {
	return &schema.Schema{
		Description: "Any of elements are in a list predicate",
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
