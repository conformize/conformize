// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package list

import (
	"fmt"
	"time"

	"github.com/conformize/conformize/common/typed"
	"github.com/conformize/conformize/internal/providers/api/schema"
	"github.com/conformize/conformize/internal/providers/api/schema/attributes"
	"github.com/spaolacci/murmur3"
)

type ListIsSubsetPredicate struct {
	Args typed.Valuable
}

func (listIsSubsetPrd *ListIsSubsetPredicate) Test(value typed.Valuable) (bool, error) {
	if value == nil || listIsSubsetPrd.Args == nil {
		return false, nil
	}

	listVal, ok := value.(*typed.ListValue)
	if !ok {
		return false, fmt.Errorf("expected a list value, got %s", value.Type().Name())
	}

	oListVal, ok := listIsSubsetPrd.Args.(*typed.ListValue)
	if !ok {
		return false, fmt.Errorf("expected a list value as argument, got %s", listIsSubsetPrd.Args.Type().Name())
	}

	if listVal.ElementsType.Hint() != oListVal.ElementsType.Hint() {
		return false, fmt.Errorf("cannot compare lists with elements of different type: %s vs %s", listVal.ElementsType.Name(), oListVal.ElementsType.Name())
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

	for _, elem := range listVal.Elements {
		if elem.Type().Hint() != listVal.ElementsType.Hint() {
			return false, fmt.Errorf(
				"element type mismatch, expected %s, got: %s",
				listVal.ElementsType.Name(), elem.Type().Name(),
			)
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

func (listIsSubsetPrd *ListIsSubsetPredicate) ArgumentsCount() int {
	return 1
}

func (listIsSubsetPrd *ListIsSubsetPredicate) Arguments(args typed.Valuable) {
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
			"Arguments": &attributes.ListAttribute{
				Required:     true,
				ElementsType: &typed.GenericTyped{},
			},
		},
	}
}
