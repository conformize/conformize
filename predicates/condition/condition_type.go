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
	EQUAL ConditionType = iota
	NOT_EQUAL
	GREATER_THAN
	LESS_THAN
	GREATER_THAN_OR_EQUAL
	LESS_THAN_OR_EQUAL
	CONTAINS
	NOT_CONTAINS
	IS_TRUE
	IS_FALSE
	MATCHES_EXPRESSION
	WITHIN_RANGE
	IS_EMPTY
	NOT_EMPTY
	IS_SUBSET
	NO_DUPLICATES
	CONTAINS_ANY
	CONTAINS_ALL
	DATE_BEFORE
	DATE_AFTER
	DATE_UP_TO
	DATE_FROM
	VALID_DATE
	SAME_DATE
	NOT_SAME_DATE
	DATE_WITHIN_INTERVAL
	DATE_IN_FUTURE
	UNKNOWN
)

var conditionTypeMap = map[string]ConditionType{}

func init() {
	for i := range int(UNKNOWN) + 1 {
		c := ConditionType(i)
		conditionTypeMap[toCamelCase(c.String())] = c
	}
}

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
