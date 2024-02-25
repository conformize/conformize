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
	"github.com/spaolacci/murmur3"
)

type IsSubsetPredicate struct{}

func (isSubsetPrd *IsSubsetPredicate) Test(value typed.Valuable, args typed.Valuable) (bool, error) {
	if value == nil || args == nil {
		return false, nil
	}

	elemListVal, ok := value.(*typed.ListValue)
	if !ok {
		return false, fmt.Errorf("expected a list value, got %s", value.Type().Name())
	}

	oElemListVal, ok := args.(*typed.ListValue)
	if !ok {
		return false, fmt.Errorf("expected a list value as argument")
	}

	if elemListVal.ElementsType.Hint().TypeHint() != oElemListVal.ElementsType.Hint().TypeHint() {
		return false, fmt.Errorf("cannot compare lists with elements of different type")
	}

	if elemListVal.Length() == 0 && oElemListVal.Length() == 0 {
		return true, nil
	}

	if elemListVal.Length() > oElemListVal.Length() {
		return false, nil
	}

	seen := make(map[string]int)
	seed := uint32(time.Now().UnixNano())
	hasher := murmur3.New128WithSeed(seed)
	for _, elem := range oElemListVal.Items() {
		elemKey, err := getElementKey(elem, hasher, seed)
		if err != nil {
			return false, err
		}
		seen[string(elemKey)]++
	}

	for _, elem := range elemListVal.Items() {
		elemKey, err := getElementKey(elem, hasher, seed)
		if err != nil {
			return false, err
		}
		keyStr := string(elemKey)
		occurences, ok := seen[keyStr]
		if !ok {
			return false, nil
		}
		occurences--
		if occurences < 0 {
			return false, nil
		}
		seen[keyStr] = occurences
	}
	return true, nil
}

func (isSubsetPrd *IsSubsetPredicate) Schema() schema.Schemable {
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
