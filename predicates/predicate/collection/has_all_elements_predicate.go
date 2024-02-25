// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package collection

import (
	"fmt"
	"time"

	"github.com/conformize/conformize/common/typed"
	"github.com/conformize/conformize/internal/providers/api/schema"
	"github.com/conformize/conformize/internal/providers/api/schema/attributes"
	"github.com/conformize/conformize/predicates"
	"github.com/spaolacci/murmur3"
)

type HasAllElementsPredicate struct {
	PredicateBuilder predicates.PredicateBuilder
}

func (hasAllElementsPrd *HasAllElementsPredicate) Test(value typed.Valuable, args typed.Valuable) (bool, error) {
	elementsVal, ok := value.(typed.Elementable)
	if !ok {
		return false, fmt.Errorf("invalid value type: expected Elementable, got %T", value)
	}

	expected, ok := args.(typed.Elementable)
	if !ok {
		return false, fmt.Errorf("invalid arguments type: expected Elementable, got %T", args)
	}
	remainingMatches := make(map[[16]byte]struct{})
	seed := uint32(time.Now().UnixNano())
	hasher := murmur3.New128WithSeed(seed)

	for _, elem := range expected.Items() {
		elemKey, err := elementHash(elem, hasher, seed)
		if err != nil {
			return false, err
		}
		remainingMatches[elemKey] = struct{}{}
	}

	elements := elementsVal.Items()
	startIdx, endIdx := 0, elementsVal.Length()-1

	for startIdx <= endIdx {
		elemHash, err := elementHash(elements[startIdx], hasher, seed)
		if err != nil {
			return false, err
		}
		if _, ok := remainingMatches[elemHash]; ok {
			delete(remainingMatches, elemHash)
			if len(remainingMatches) == 0 {
				return true, nil
			}
		}

		if startIdx == endIdx {
			break
		}

		elemHash, err = elementHash(elements[endIdx], hasher, seed)
		if err != nil {
			return false, err
		}

		if _, ok := remainingMatches[elemHash]; ok {
			delete(remainingMatches, elemHash)
			if len(remainingMatches) == 0 {
				return true, nil
			}
		}

		startIdx++
		endIdx--
	}

	return len(remainingMatches) == 0, nil
}

func (hasAllElementsPrd *HasAllElementsPredicate) Schema() schema.Schemable {
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
