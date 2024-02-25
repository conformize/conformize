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

func TestNewObject(t *testing.T) {
	testCases := []struct {
		valueType *typed.ObjectTyped
		value     any
		expected  typed.Valuable
	}{
		{
			valueType: &typed.ObjectTyped{
				FieldsTypes: map[string]typed.Typeable{
					"name": &typed.StringTyped{}, "age": &typed.NumberTyped{},
				}},
			value: struct {
				Name string `cnfrmz:"name"`
				Age  int    `cnfrmz:"age"`
			}{"John Doe", 42},
			expected: typed.NewObjectValue(map[string]typed.Valuable{
				"name": value(typed.NewStringValue("John Doe")),
				"age":  value(typed.NewNumberValue(42)),
			},
				map[string]typed.Typeable{"name": &typed.StringTyped{}, "age": &typed.NumberTyped{}}),
		},
		{
			valueType: &typed.ObjectTyped{
				FieldsTypes: map[string]typed.Typeable{
					"sandbox":     &typed.BooleanTyped{},
					"environment": &typed.StringTyped{},
					"host":        &typed.StringTyped{},
				},
			},
			value: struct {
				Sandbox     bool   `cnfrmz:"sandbox"`
				Environment string `cnfrmz:"environment"`
				Host        string `cnfrmz:"host"`
			}{true, "dev", "localhost"},
			expected: typed.NewObjectValue(map[string]typed.Valuable{
				"sandbox":     value(typed.NewBooleanValue(true)),
				"environment": value(typed.NewStringValue("dev")),
				"host":        value(typed.NewStringValue("localhost")),
			},
				map[string]typed.Typeable{
					"sandbox":     &typed.BooleanTyped{},
					"environment": &typed.StringTyped{},
					"host":        &typed.StringTyped{},
				}),
		},
	}

	for _, testCase := range testCases {
		if val, err := Object(reflect.ValueOf(testCase.value), testCase.valueType); err != nil || !reflect.DeepEqual(val, testCase.expected) {
			t.Fail()
		}
	}
}

type TestStructA struct {
	Sandbox     bool   `cnfrmz:"sandbox"`
	Environment string `cnfrmz:"environment"`
	Host        string `cnfrmz:"host"`
	Nodes       int16  `cnfrmz:"nodes"`
}

func TestObjectToStruct(t *testing.T) {
	testCases := []struct {
		objectValue typed.ObjectValue
		value       any
		expected    any
	}{
		{
			objectValue: typed.ObjectValue{
				Fields: map[string]typed.Valuable{
					"sandbox":     value(typed.NewBooleanValue(true)),
					"environment": value(typed.NewStringValue("dev")),
					"host":        value(typed.NewStringValue("localhost")),
					"nodes":       value(typed.NewNumberValue(3)),
				},
				FieldsTypes: map[string]typed.Typeable{
					"sandbox":     &typed.BooleanTyped{},
					"environment": &typed.StringTyped{},
					"host":        &typed.StringTyped{},
					"nodes":       &typed.NumberTyped{},
				},
			},
			value:    &TestStructA{},
			expected: &TestStructA{true, "dev", "localhost", 3},
		},
	}

	for _, testCase := range testCases {
		if err := testCase.objectValue.As(testCase.value); err != nil || !reflect.DeepEqual(testCase.value, testCase.expected) {
			t.Fail()
		}
	}
}
