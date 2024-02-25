// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

//go:generate stringer -type=ConditionType -output=condition_type_string.go
package condition

import (
	"strings"
	"unicode"
)

type ConditionType int

const (
	EQ ConditionType = iota
	NOT
	GT
	LT
	GTE
	LTE
	HAS
	LACKS
	TRUE
	FALSE
	MATCHES
	RANGE
	EMPTY
	NOT_EMPTY
	SUBSET_OF
	UNIQUE
	HAS_ANY
	BEFORE
	AFTER
	UNTIL
	SINCE
	VALID
	SAME
	DIFFERENT
	WITHIN
	FUTURE
	UNKNOWN
)

var conditionTypeMap = func() map[string]ConditionType {
	m := make(map[string]ConditionType, int(UNKNOWN)+1)
	for i := ConditionType(0); i <= UNKNOWN; i++ {
		m[toCamelCase(i.String())] = i
	}
	return m
}()

func toCamelCase(s string) string {
	var result strings.Builder
	nextUpper := false
	for i, r := range s {
		if i == 0 {
			result.WriteRune(unicode.ToLower(r))
			continue
		}

		if r == '_' || r == ' ' {
			nextUpper = true
			continue
		}

		if nextUpper {
			result.WriteRune(unicode.ToUpper(r))
			nextUpper = false
		} else {
			result.WriteRune(unicode.ToLower(r))
		}
	}
	return result.String()
}

func FromString(s string) ConditionType {
	if c, ok := conditionTypeMap[s]; ok {
		return c
	}
	return UNKNOWN
}
