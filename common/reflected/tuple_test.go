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
)

func TestTupleValue(t *testing.T) {
	testCases := []struct {
		valueType typed.Typeable
		value     interface{}
		expected  typed.Valuable
	}{
		{
			valueType: &typed.TupleTyped{ElementsTypes: []typed.Typeable{&typed.StringTyped{}, &typed.NumberTyped{}, &typed.BooleanTyped{}}},
			value:     []interface{}{"hello", 42, true},
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
		if val, err := Tuple(tc.value, tc.valueType); err != nil || !reflect.DeepEqual(val, tc.expected) {
			t.Fail()
		}
	}
}

func TestTupleValueToSlice(t *testing.T) {
	testCases := []struct {
		value    interface{}
		tuple    typed.Valuable
		expected interface{}
	}{
		{
			value: []interface{}{},
			tuple: &typed.TupleValue{
				ElementsTypes: []typed.Typeable{&typed.NumberTyped{}, &typed.StringTyped{}},
				Elements: []typed.Valuable{
					value(typed.NewNumberValue(42)),
					value(typed.NewStringValue("test")),
				},
			},
			expected: []interface{}{42, "test"},
		},
		{
			value: []interface{}{32, "hello"},
			tuple: &typed.TupleValue{
				ElementsTypes: []typed.Typeable{&typed.NumberTyped{}, &typed.StringTyped{}},
				Elements: []typed.Valuable{
					value(typed.NewNumberValue(42)),
					value(typed.NewStringValue("test")),
				},
			},
			expected: []interface{}{42, "test"},
		},
	}

	for _, tc := range testCases {
		if err := tc.tuple.As(&tc.value); err != nil || !reflect.DeepEqual(tc.value, tc.expected) {
			t.Fail()
		}
	}
}
