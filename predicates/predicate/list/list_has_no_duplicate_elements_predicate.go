// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package list

import (
	"encoding/binary"
	"fmt"
	"math"
	"time"

	"github.com/conformize/conformize/common/typed"
	"github.com/conformize/conformize/internal/providers/api/schema"
	"github.com/conformize/conformize/internal/providers/api/schema/attributes"
	"github.com/conformize/conformize/predicates/predicate"
	"github.com/spaolacci/murmur3"
)

type ListHasNoDuplicateElementsPredicate struct {
	predicate.PredicateArgumentsValidator
}

func (listHasNoDupPrd *ListHasNoDuplicateElementsPredicate) Test(value typed.Valuable, args *typed.TupleValue) (bool, error) {
	valid, validErr := listHasNoDupPrd.Validate(value, args, listHasNoDupPrd.Schema())
	if !valid {
		return valid, validErr
	}

	listVal := value.(*typed.ListValue)
	seen := make(map[[16]byte]struct{})
	seed := uint32(time.Now().UnixNano())
	hasher := murmur3.New128WithSeed(seed)
	for startIdx, endIdx := 0, len(listVal.Elements)-1; endIdx >= startIdx; startIdx, endIdx = startIdx+1, endIdx-1 {
		for _, idx := range []int{startIdx, endIdx} {
			elemKey, err := elementHash(listVal.Elements[idx], hasher, seed)
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
		}
	}
	return true, nil
}

func (listHasNoDupPrd *ListHasNoDuplicateElementsPredicate) Arguments() int {
	return 0
}

func (listHasNoDupPrd *ListHasNoDuplicateElementsPredicate) Schema() schema.Schemable {
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
