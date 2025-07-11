// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package collection

import (
	"encoding/binary"
	"fmt"
	"math"
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
	seen := make(map[[16]byte]struct{})
	seed := uint32(time.Now().UnixNano())
	hasher := murmur3.New128WithSeed(seed)

	startIdx, endIdx := 0, listVal.Length()-1
	for startIdx <= endIdx {
		elemKey, err := elementHash(listVal.Elements[startIdx], hasher, seed)
		if err != nil {
			return false, err
		}
		if _, found := seen[elemKey]; found {
			return false, nil
		}

		seen[elemKey] = struct{}{}
		if startIdx == endIdx {
			break
		}

		elemKey, err = elementHash(listVal.Elements[endIdx], hasher, seed)
		if err != nil {
			return false, err
		}
		if _, found := seen[elemKey]; found {
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

func elementHash(val typed.Valuable, hasher murmur3.Hash128, seed uint32) ([16]byte, error) {
	hasher.Reset()
	switch v := val.(type) {
	case *typed.BooleanValue:
		var boolVal bool
		v.As(&boolVal)
		hasher.Write([]byte(fmt.Sprintf("%t", boolVal)))
	case *typed.NumberValue:
		var numVal float64
		v.As(&numVal)
		numBytes := make([]byte, 8)
		binary.LittleEndian.PutUint64(numBytes, math.Float64bits(numVal))
		hasher.Write(numBytes)
	case *typed.StringValue:
		var strVal string
		v.As(&strVal)
		hasher.Write([]byte(strVal))
	case *typed.ListValue:
		listHasher := murmur3.New128WithSeed(seed)
		for _, elem := range v.Elements {
			elemHash, err := elementHash(elem, hasher, seed)
			if err != nil {
				return [16]byte{}, err
			}
			listHasher.Write(elemHash[:])
		}
		var hash [16]byte
		listHasher.Sum(hash[:0])
		return hash, nil
	default:
		return [16]byte{}, fmt.Errorf("unsupported type: %v", val.Type().Name())
	}
	var hash [16]byte
	hasher.Sum(hash[:0])
	return hash, nil
}
