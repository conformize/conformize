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

func value(v typed.Valuable, _ error) typed.Valuable {
	return v
}

func TestListValue(t *testing.T) {
	testCases := []struct {
		valueType typed.Typeable
		value     interface{}
		expected  typed.Valuable
	}{
		{
			valueType: &typed.ListTyped{ElementsType: &typed.NumberTyped{}},
			value:     []int64{1, 2, 3},
			expected: &typed.ListValue{
				ElementsType: &typed.NumberTyped{},
				Elements: []typed.Valuable{
					value(typed.NewNumberValue(int64(1))),
					value(typed.NewNumberValue(int64(2))),
					value(typed.NewNumberValue(int64(3))),
				},
			},
		},
		{
			valueType: &typed.ListTyped{ElementsType: &typed.StringTyped{}},
			value:     []string{"hello", "world"},
			expected: &typed.ListValue{
				ElementsType: &typed.StringTyped{},
				Elements: []typed.Valuable{
					value(typed.NewStringValue("hello")),
					value(typed.NewStringValue("world")),
				},
			},
		},
		{
			valueType: &typed.ListTyped{ElementsType: &typed.BooleanTyped{}},
			value:     []bool{true, false},
			expected: &typed.ListValue{
				ElementsType: &typed.BooleanTyped{},
				Elements: []typed.Valuable{
					value(typed.NewBooleanValue(true)),
					value(typed.NewBooleanValue(false)),
				},
			},
		},
		{
			valueType: &typed.ListTyped{ElementsType: &typed.NumberTyped{}},
			value:     []float64{3.14, 2.71},
			expected: &typed.ListValue{
				ElementsType: &typed.NumberTyped{},
				Elements: []typed.Valuable{
					value(typed.NewNumberValue(3.14)),
					value(typed.NewNumberValue(float64(2.71)))},
			},
		},
	}

	for _, tc := range testCases {
		if val, err := List(reflect.ValueOf(tc.value), tc.valueType); err != nil || !reflect.DeepEqual(val, tc.expected) {
			t.Fail()
		}
	}
}

type TestStruct struct {
	Sandbox     bool   `cnfrmz:"sandbox"`
	Environment string `cnfrmz:"environment"`
	Host        string `cnfrmz:"host"`
}

func TestListValueToSlice(t *testing.T) {
	testCases := []struct {
		name     string
		value    interface{}
		list     typed.Valuable
		expected interface{}
	}{
		{
			name:  "Test reflect list value to []int64",
			value: []int64{},
			list: &typed.ListValue{
				ElementsType: &typed.NumberTyped{},
				Elements: []typed.Valuable{
					value(typed.NewNumberValue(int64(1))),
					value(typed.NewNumberValue(int64(2))),
					value(typed.NewNumberValue(int64(3))),
				},
			},
			expected: []int64{1, 2, 3},
		},
		{
			name:  "Test reflect list value to []string",
			value: []string{},
			list: &typed.ListValue{
				ElementsType: &typed.StringTyped{},
				Elements: []typed.Valuable{
					value(typed.NewStringValue("hello")),
					value(typed.NewStringValue("world")),
				},
			},
			expected: []string{"hello", "world"},
		},
		{
			name:  "Test reflect list value to []bool",
			value: []bool{},
			list: &typed.ListValue{
				ElementsType: &typed.BooleanTyped{},
				Elements: []typed.Valuable{
					value(typed.NewBooleanValue(true)),
					value(typed.NewBooleanValue(false)),
				},
			},
			expected: []bool{true, false},
		},
		{
			name:  "Test reflect list value to []float64",
			value: []float64{},
			list: &typed.ListValue{
				ElementsType: &typed.NumberTyped{},
				Elements: []typed.Valuable{
					value(typed.NewNumberValue(3.14)),
					value(typed.NewNumberValue(float64(2.71))),
				},
			},
			expected: []float64{3.14, 2.71},
		},
		{
			name:  "Test reflect list value to []struct{}",
			value: []TestStruct{},
			list: &typed.ListValue{
				ElementsType: &typed.ObjectTyped{
					FieldsTypes: map[string]typed.Typeable{
						"sandbox":     &typed.BooleanTyped{},
						"environment": &typed.StringTyped{},
						"host":        &typed.StringTyped{},
					},
				},
				Elements: []typed.Valuable{
					typed.NewObjectValue(map[string]typed.Valuable{
						"sandbox":     value(typed.NewBooleanValue(true)),
						"environment": value(typed.NewStringValue("dev")),
						"host":        value(typed.NewStringValue("localhost")),
					},
						map[string]typed.Typeable{
							"sandbox":     &typed.BooleanTyped{},
							"environment": &typed.StringTyped{},
							"host":        &typed.StringTyped{},
						}),
					typed.NewObjectValue(map[string]typed.Valuable{
						"sandbox":     value(typed.NewBooleanValue(false)),
						"environment": value(typed.NewStringValue("prod")),
						"host":        value(typed.NewStringValue("127.0.0.1")),
					}, map[string]typed.Typeable{
						"sandbox":     &typed.BooleanTyped{},
						"environment": &typed.StringTyped{},
						"host":        &typed.StringTyped{},
					}),
				},
			},
			expected: []TestStruct{
				{true, "dev", "localhost"},
				{false, "prod", "127.0.0.1"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.list.As(&tc.value); err != nil || !reflect.DeepEqual(tc.value, tc.expected) {
				t.Fail()
			}
		})
	}
}

func TestListValueToNestedSlice(t *testing.T) {
	var slice [][]float64
	expected := [][]float64{{1, 2, 3}, {4, 5, 6}}
	value := &typed.ListValue{
		ElementsType: &typed.ListTyped{
			ElementsType: &typed.NumberTyped{},
		},
		Elements: []typed.Valuable{
			&typed.ListValue{
				ElementsType: &typed.NumberTyped{},
				Elements: []typed.Valuable{
					value(typed.NewNumberValue(int64(1))),
					value(typed.NewNumberValue(int64(2))),
					value(typed.NewNumberValue(int64(3))),
				},
			},
			&typed.ListValue{
				ElementsType: &typed.NumberTyped{},
				Elements: []typed.Valuable{
					value(typed.NewNumberValue(int64(4))),
					value(typed.NewNumberValue(int64(5))),
					value(typed.NewNumberValue(int64(6))),
				},
			},
		},
	}
	if err := value.As(&slice); err != nil || !reflect.DeepEqual(slice, expected) {
		t.Fail()
	}
}
