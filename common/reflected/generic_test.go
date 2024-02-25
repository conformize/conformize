// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package reflected

import (
	"reflect"
	"testing"

	"github.com/conformize/conformize/common/typed"
	"github.com/google/go-cmp/cmp"
)

func TestGenericValue(t *testing.T) {
	testCases := []struct {
		value         any
		resolvedValue typed.Valuable
		resolvedRaw   any
		resolved      typed.Valuable
		expected      typed.Valuable
	}{
		{
			value:         []string{"hello", "there,", "wordld!"},
			resolvedValue: &typed.ListValue{ElementsType: &typed.StringTyped{}},
			resolvedRaw:   []any{},
			expected: &typed.ListValue{
				ElementsType: &typed.StringTyped{},
				Elements: []typed.Valuable{
					value(typed.NewStringValue("hello")),
					value(typed.NewStringValue("there,")),
					value(typed.NewStringValue("wordld!")),
				},
			},
		},
		{
			value:         []int64{1, 1, 1970},
			resolvedValue: &typed.ListValue{ElementsType: &typed.NumberTyped{}},
			resolvedRaw:   []int64{},
			expected: &typed.ListValue{
				ElementsType: &typed.NumberTyped{},
				Elements: []typed.Valuable{
					value(typed.NewNumberValue(1)),
					value(typed.NewNumberValue(1)),
					value(typed.NewNumberValue(1970)),
				},
			},
		},
	}

	for _, tc := range testCases {
		val, err := Generic(reflect.ValueOf(tc.value))
		if err != nil {
			t.Fail()
		}
		err = val.As(tc.resolvedValue)
		if err != nil {
			t.Fail()
		}
		if !reflect.DeepEqual(tc.resolvedValue, tc.expected) {
			t.Fail()
		}

		err = val.As(&tc.resolvedRaw)
		if err != nil {
			t.Fail()
		}
	}
}

func TestAssignGenericValue(t *testing.T) {
	testCases := []struct {
		name        string
		value       typed.Valuable
		resolvedRaw any
		expected    any
	}{
		{
			name:        "assign list of strings to generic value succeeds",
			expected:    []string{"hello", "there,", "wordld!"},
			resolvedRaw: []string{},
			value: &typed.ListValue{
				ElementsType: &typed.StringTyped{},
				Elements: []typed.Valuable{
					value(typed.NewStringValue("hello")),
					value(typed.NewStringValue("there,")),
					value(typed.NewStringValue("wordld!")),
				},
			},
		},
		{
			name:        "assign list of int64 to generic value succeeds",
			expected:    []int64{1, 1, 1970},
			resolvedRaw: []int64{},
			value: &typed.ListValue{
				ElementsType: &typed.NumberTyped{},
				Elements: []typed.Valuable{
					value(typed.NewNumberValue(1)),
					value(typed.NewNumberValue(1)),
					value(typed.NewNumberValue(1970)),
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			val := &typed.GenericValue{}
			err := val.Assign(tc.value)
			if err != nil {
				t.Fail()
			}

			err = val.As(&tc.resolvedRaw)
			if err != nil {
				t.Fail()
			}

			if !cmp.Equal(tc.resolvedRaw, tc.expected) {
				t.Fail()
			}
		})
	}
}
