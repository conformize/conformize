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
		value     interface{}
		expected  typed.Valuable
	}{
		{
			valueType: &typed.ObjectTyped{
				FieldsTypes: map[string]typed.Typeable{
					"Name": &typed.StringTyped{}, "Age": &typed.NumberTyped{},
				}},
			value: struct {
				Name string
				Age  int
			}{"John Doe", 42},
			expected: typed.NewObjectValue(map[string]typed.Valuable{
				"Name": value(typed.NewStringValue("John Doe")),
				"Age":  value(typed.NewNumberValue(42)),
			},
				map[string]typed.Typeable{"Name": &typed.StringTyped{}, "Age": &typed.NumberTyped{}}),
		},
		{
			valueType: &typed.ObjectTyped{
				FieldsTypes: map[string]typed.Typeable{
					"Sandbox":     &typed.BooleanTyped{},
					"Environment": &typed.StringTyped{},
					"Host":        &typed.StringTyped{},
				},
			},
			value: struct {
				Sandbox     bool
				Environment string
				Host        string
			}{true, "dev", "localhost"},
			expected: typed.NewObjectValue(map[string]typed.Valuable{
				"Sandbox":     value(typed.NewBooleanValue(true)),
				"Environment": value(typed.NewStringValue("dev")),
				"Host":        value(typed.NewStringValue("localhost")),
			},
				map[string]typed.Typeable{
					"Sandbox":     &typed.BooleanTyped{},
					"Environment": &typed.StringTyped{},
					"Host":        &typed.StringTyped{},
				}),
		},
	}

	for _, testCase := range testCases {
		if val, err := Object(testCase.value, testCase.valueType); err != nil || !reflect.DeepEqual(val, testCase.expected) {
			t.Fail()
		}
	}
}

type TestStructA struct {
	Sandbox     bool
	Environment string
	Host        string
	Nodes       int16
}

type NestedStruct struct {
	X      int
	Nested TestStructB
}

type TestStructB struct {
	Y TestStructA
}

func TestObjectsToStruct(t *testing.T) {
	testCases := []struct {
		objectValue typed.ObjectValue
		value       interface{}
		expected    interface{}
	}{
		{
			objectValue: typed.ObjectValue{
				Fields: map[string]typed.Valuable{
					"Sandbox":     value(typed.NewBooleanValue(true)),
					"Environment": value(typed.NewStringValue("dev")),
					"Host":        value(typed.NewStringValue("localhost")),
					"Nodes":       value(typed.NewNumberValue(3)),
				},
				FieldsTypes: map[string]typed.Typeable{
					"Sandbox":     &typed.BooleanTyped{},
					"Environment": &typed.StringTyped{},
					"Host":        &typed.StringTyped{},
					"Nodes":       &typed.NumberTyped{},
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

func TestObjectToStruct(t *testing.T) {
	var st TestStructA
	expected := TestStructA{true, "dev", "localhost", 3}
	objectValue := typed.ObjectValue{
		Fields: map[string]typed.Valuable{
			"Sandbox":     value(typed.NewBooleanValue(true)),
			"Environment": value(typed.NewStringValue("dev")),
			"Host":        value(typed.NewStringValue("localhost")),
			"Nodes":       value(typed.NewNumberValue(3)),
		},
		FieldsTypes: map[string]typed.Typeable{
			"Sandbox":     &typed.BooleanTyped{},
			"Environment": &typed.StringTyped{},
			"Host":        &typed.StringTyped{},
			"Nodes":       &typed.NumberTyped{},
		}}
	objectValue.As(&st)
	if !reflect.DeepEqual(st, expected) {
		t.Fail()
	}
}
