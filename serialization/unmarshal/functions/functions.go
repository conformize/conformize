// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package functions

import (
	"fmt"
	"strconv"

	"github.com/conformize/conformize/common/ds"
)

func UnmarshalValue[K comparable, V any](nodeRef *ds.Node[K, V], value any) {
	if val, ok := value.(map[K]any); !ok {
		nodeRef.Value = value.(V)
	} else {
		for key, v := range val {
			var childNodeRef = nodeRef.AddChild(key)
			UnmarshalValue(childNodeRef, v)
		}
	}
}

func DecodeStringValue(value string) (any, error) {
	if val, err := strconv.ParseFloat(value, 64); err == nil {
		return val, nil
	}
	if val, err := strconv.ParseBool(value); err == nil {
		return val, nil
	}
	return nil, fmt.Errorf("couldn't decode string value: %s", value)
}

func IsWhiteSpace(value string) bool {
	for _, char := range value {
		if char != ' ' && char != '\t' && char != '\n' && char != '\r' {
			return false
		}
	}
	return true
}
