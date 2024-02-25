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

type HasNoDuplicateElementsPredicate struct{}

func (hasNoDupPrd *HasNoDuplicateElementsPredicate) Test(value typed.Valuable, Args typed.Valuable) (bool, error) {
	listVal, ok := value.(*typed.ListValue)
	if !ok {
		return false, fmt.Errorf("expected a list value, got %T", value)
	}
	seen := make(map[string]struct{})
	seed := uint32(time.Now().UnixNano())
	hasher := murmur3.New128WithSeed(seed)

	startIdx, endIdx := 0, listVal.Length()-1
	for startIdx <= endIdx {
		elemKey, err := getElementKey(listVal.Elements[startIdx], hasher, seed)
		if err != nil {
			return false, err
		}
		keyStr := string(elemKey)
		if _, found := seen[keyStr]; found {
			return false, nil
		}

		seen[keyStr] = struct{}{}
		if startIdx == endIdx {
			break
		}

		elemKey, err = getElementKey(listVal.Elements[endIdx], hasher, seed)
		if err != nil {
			return false, err
		}
		keyStr = string(elemKey)
		if _, found := seen[keyStr]; found {
			return false, nil
		}

		startIdx++
		endIdx--
	}
	return true, nil
}

func (hasNoDupPrd *HasNoDuplicateElementsPredicate) Schema() schema.Schemable {
	return &schema.Schema{
		Description: "All elements in a list are unique predicate",
		Version:     1,
		Attributes: map[string]schema.Attributeable{
			"Value": &attributes.ListAttribute{
				Required:     true,
				ElementsType: &typed.GenericTyped{},
			},
		},
	}
}

