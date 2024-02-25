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

func TestTupleValue(t *testing.T) {
	testCases := []struct {
		valueType typed.Typeable
		value     any
		expected  typed.Valuable
	}{
		{
			valueType: &typed.TupleTyped{ElementsTypes: []typed.Typeable{&typed.StringTyped{}, &typed.NumberTyped{}, &typed.BooleanTyped{}}},
			value:     []any{"hello", 42, true},
			expected: &typed.TupleValue{
				ElementsTypes: []typed.Typeable{&typed.StringTyped{}, &typed.NumberTyped{}, &typed.BooleanTyped{}},
				Elements: []typed.Valuable{
					value(typed.NewStringValue("hello")),
					value(typed.NewNumberValue(42)),
					value(typed.NewBooleanValue(true)),
				},
			},
		},
	}

	for _, tc := range testCases {
		if val, err := Tuple(reflect.ValueOf(tc.value), tc.valueType); err != nil || !reflect.DeepEqual(val, tc.expected) {
			t.Fail()
		}
	}
}

func TestTupleValueToSlice(t *testing.T) {
	testCases := []struct {
		name     string
		value    any
		tuple    typed.Valuable
		expected any
	}{
		{
			name:  "empty tuple to list",
			value: []any{},
			tuple: &typed.TupleValue{
				ElementsTypes: []typed.Typeable{&typed.NumberTyped{}, &typed.StringTyped{}},
				Elements: []typed.Valuable{
					value(typed.NewNumberValue(42)),
					value(typed.NewStringValue("test")),
				},
			},
			expected: []any{float64(42), "test"},
		},
		{
			name:  "tuple with mixed types to list",
			value: []any{32, "hello"},
			tuple: &typed.TupleValue{
				ElementsTypes: []typed.Typeable{&typed.NumberTyped{}, &typed.StringTyped{}},
				Elements: []typed.Valuable{
					value(typed.NewNumberValue(42)),
					value(typed.NewStringValue("test")),
				},
			},
			expected: []any{42, "test"},
		},
		{
			name:  "tuple with nested list",
			value: []any{},
			tuple: &typed.TupleValue{
				ElementsTypes: []typed.Typeable{
					&typed.ListTyped{ElementsType: &typed.NumberTyped{}},
					&typed.StringTyped{},
				},
				Elements: []typed.Valuable{
					&typed.ListValue{
						ElementsType: &typed.NumberTyped{},
						Elements: []typed.Valuable{
							value(typed.NewNumberValue(1)),
							value(typed.NewNumberValue(2)),
						},
					},
					value(typed.NewStringValue("ok")),
				},
			},
			expected: []any{
				[]float64{1, 2},
				"ok",
			},
		},
		{
			name:  "tuple with empty nested list",
			value: []any{},
			tuple: &typed.TupleValue{
				ElementsTypes: []typed.Typeable{
					&typed.ListTyped{ElementsType: &typed.NumberTyped{}},
					&typed.StringTyped{},
				},
				Elements: []typed.Valuable{
					&typed.ListValue{
						ElementsType: &typed.NumberTyped{},
						Elements:     []typed.Valuable{value(typed.NewNumberValue(1)), value(typed.NewNumberValue(2)), value(typed.NewNumberValue(3))},
					},
					value(typed.NewStringValue("empty")),
				},
			},
			expected: []any{
				[]float64{1, 2, 3},
				"empty",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.tuple.As(&tc.value); err != nil {
				t.Errorf("test %s failed: As() returned error: %v", tc.name, err)
			}

			if !cmp.Equal(tc.value, tc.expected) {
				t.Errorf("test %s failed: expected %v, got %v", tc.name, tc.expected, tc.value)
			}
		})
	}
}
