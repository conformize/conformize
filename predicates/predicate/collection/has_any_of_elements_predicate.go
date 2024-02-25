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

type HasAnyOfElementsPredicate struct {
	PredicateBuilder predicates.PredicateBuilder
}

func (hasAnyOfElemPrd *HasAnyOfElementsPredicate) Test(value typed.Valuable, args typed.Valuable) (bool, error) {
	listVal, ok := value.(*typed.ListValue)
	if !ok {
		return false, fmt.Errorf("expected a list value, got %T", value)
	}
	elementsListVal, ok := args.(*typed.ListValue)
	if !ok {
		return false, fmt.Errorf("expected a list value as argument, got %T", args)
	}

	if listVal.ElementsType.Hint().TypeHint() != elementsListVal.ElementsType.Hint().TypeHint() {
		return false, fmt.Errorf("cannot compare lists with elements of different type")
	}

	seen := make(map[string]struct{})
	seed := uint32(time.Now().UnixNano())
	hasher := murmur3.New128WithSeed(seed)
	for _, elem := range elementsListVal.Elements {
		elemKey, err := getElementKey(elem, hasher, seed)
		if err != nil {
			return false, err
		}
		seen[string(elemKey)] = struct{}{}
	}

	for _, elem := range listVal.Elements {
		elemKey, err := getElementKey(elem, hasher, seed)
		if err != nil {
			return false, err
		}
		if _, ok := seen[string(elemKey)]; ok {
			return true, nil
		}
	}
	return false, nil
}

func (hasAnyOfElemPrd *HasAnyOfElementsPredicate) Schema() schema.Schemable {
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
