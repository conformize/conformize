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

	"github.com/conformize/conformize/common/typed"
	"github.com/spaolacci/murmur3"
)

// getElementKey returns a []byte key for any element type
// Primitives get direct byte representation, complex types get hashed
func getElementKey(val typed.Valuable, hasher murmur3.Hash128, seed uint32) ([]byte, error) {
	switch v := val.(type) {
	case *typed.BooleanValue:
		var boolVal bool
		v.As(&boolVal)
		if boolVal {
			return []byte{byte(typed.Boolean), 1}, nil
		}
		return []byte{byte(typed.Boolean), 0}, nil

	case *typed.NumberValue:
		var numVal float64
		v.As(&numVal)
		key := make([]byte, 9) // 1 byte type + 8 bytes number
		key[0] = byte(typed.Number)
		binary.LittleEndian.PutUint64(key[1:], math.Float64bits(numVal))
		return key, nil

	case *typed.StringValue:
		var strVal string
		v.As(&strVal)
		return []byte(strVal), nil

	default:
		// Complex types: use hash with type prefix
		hash, err := elementHash(val, hasher, seed)
		if err != nil {
			return nil, err
		}
		return hash[:], nil
	}
}

func elementHash(val typed.Valuable, hasher murmur3.Hash128, seed uint32) ([16]byte, error) {
	hasher.Reset()
	switch v := val.(type) {
	case *typed.ListValue:
		// For lists, hash the concatenated element keys
		for _, elem := range v.Elements {
			elemKey, err := getElementKey(elem, hasher, seed)
			if err != nil {
				return [16]byte{}, err
			}
			hasher.Write(elemKey)
		}
	case *typed.MapValue:
		// For maps, hash key-value pairs in sorted order
		// This ensures consistent hashing regardless of iteration order
		keys := make([]string, 0, len(v.Elements))
		for k := range v.Elements {
			keys = append(keys, k)
		}
		// Sort keys for deterministic hashing
		for i := 0; i < len(keys)-1; i++ {
			for j := i + 1; j < len(keys); j++ {
				if keys[i] > keys[j] {
					keys[i], keys[j] = keys[j], keys[i]
				}
			}
		}
		for _, k := range keys {
			hasher.Write([]byte(k))
			elemKey, err := getElementKey(v.Elements[k], hasher, seed)
			if err != nil {
				return [16]byte{}, err
			}
			hasher.Write(elemKey)
		}
	case *typed.ObjectValue:
		// For objects, hash field-value pairs in sorted order
		keys := make([]string, 0, len(v.Fields))
		for k := range v.Fields {
			keys = append(keys, k)
		}
		// Sort keys for deterministic hashing
		for i := 0; i < len(keys)-1; i++ {
			for j := i + 1; j < len(keys); j++ {
				if keys[i] > keys[j] {
					keys[i], keys[j] = keys[j], keys[i]
				}
			}
		}
		for _, k := range keys {
			hasher.Write([]byte(k))
			elemKey, err := getElementKey(v.Fields[k], hasher, seed)
			if err != nil {
				return [16]byte{}, err
			}
			hasher.Write(elemKey)
		}
	case *typed.TupleValue:
		// For tuples, hash elements in order
		for _, elem := range v.Elements {
			elemKey, err := getElementKey(elem, hasher, seed)
			if err != nil {
				return [16]byte{}, err
			}
			hasher.Write(elemKey)
		}
	default:
		return [16]byte{}, fmt.Errorf("unsupported complex type: %v", val.Type().Name())
	}
	var hash [16]byte
	hasher.Sum(hash[:0])
	return hash, nil
}
