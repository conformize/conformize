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

func TestNewPrimitiveFailsWithInvalidType(t *testing.T) {
	testCases := []struct {
		valueType typed.Typeable
		value     interface{}
	}{
		{
			valueType: &typed.MapTyped{ElementsType: &typed.NumberTyped{}},
			value:     int32(42),
		},
		{
			valueType: &typed.ListTyped{ElementsType: &typed.NumberTyped{}},
			value:     int32(42),
		},
	}

	for _, tc := range testCases {
		if _, err := Primitive(tc.value, tc.valueType); err == nil {
			t.Fail()
		}
	}
}

func TestNewPrimitiveSucceedsWithValidType(t *testing.T) {
	testCases := []struct {
		valueType typed.Typeable
		value     interface{}
		expected  typed.Valuable
	}{
		{
			valueType: &typed.BooleanTyped{},
			value:     true,
			expected:  value(typed.NewBooleanValue(true)),
		},
		{
			valueType: &typed.NumberTyped{},
			value:     int64(42),
			expected:  value(typed.NewNumberValue(int64(42))),
		},
		{
			valueType: &typed.NumberTyped{},
			value:     float64(3.14),
			expected:  value(typed.NewNumberValue(float64(3.14))),
		},
		{
			valueType: &typed.StringTyped{},
			value:     "hello",
			expected:  value(typed.NewStringValue("hello")),
		},
	}

	for _, tc := range testCases {
		if val, err := Primitive(tc.value, tc.valueType); err != nil || !reflect.DeepEqual(val, tc.expected) {
			t.Fail()
		}
	}
}

func TestNewPrimitiveFailsWithInvalidValueType(t *testing.T) {
	testCases := []struct {
		valueType typed.Typeable
		value     interface{}
	}{
		{
			valueType: &typed.BooleanTyped{},
			value:     int32(42),
		},
		{
			valueType: &typed.NumberTyped{},
			value:     true,
		},
		{
			valueType: &typed.NumberTyped{},
			value:     []int{1, 2, 3},
		},
		{
			valueType: &typed.StringTyped{},
			value:     map[string]int{"a": 1, "b": 2},
		},
		{
			valueType: &typed.NumberTyped{},
			value:     true,
		},
		{
			valueType: &typed.StringTyped{},
			value:     true,
		},
	}

	for _, tc := range testCases {
		if _, err := Primitive(tc.value, tc.valueType); err == nil {
			t.Fail()
		}
	}
}

func TestIntPrimitive(t *testing.T) {
	testCases := []struct {
		value    interface{}
		expected typed.Valuable
	}{
		{
			value:    int64(42),
			expected: value(typed.NewNumberValue(int64(42))),
		},
		{
			value:    int(42),
			expected: value(typed.NewNumberValue(int64(42))),
		},
		{
			value:    int8(42),
			expected: value(typed.NewNumberValue(int64(42))),
		},
		{
			value:    int16(42),
			expected: value(typed.NewNumberValue(int64(42))),
		},
		{
			value:    int32(42),
			expected: value(typed.NewNumberValue(int64(42))),
		},
	}

	for _, tc := range testCases {
		if _, err := Primitive(tc.value, &typed.NumberTyped{}); err != nil {
			t.Fail()
		}
	}
}

func TestFloatPrimitive(t *testing.T) {
	testCases := []struct {
		value    interface{}
		expected typed.Valuable
	}{
		{
			value:    float64(3.14),
			expected: value(typed.NewNumberValue(float64(3.14))),
		},
		{
			value:    float32(3.14),
			expected: value(typed.NewNumberValue(float64(3.14))),
		},
	}

	for _, tc := range testCases {
		if _, err := Primitive(tc.value, &typed.NumberTyped{}); err != nil {
			t.Fail()
		}
	}
}

func TestPrimitiveToType(t *testing.T) {
	testCases := []struct {
		primitive typed.Valuable
		dst       interface{}
		expected  interface{}
	}{
		{
			primitive: value(typed.NewBooleanValue(true)),
			dst:       bool(false),
			expected:  true,
		},
		{
			primitive: value(typed.NewNumberValue(int64(42))),
			dst:       int64(0),
			expected:  int64(42),
		},
		{
			primitive: value(typed.NewNumberValue(float64(3.14))),
			dst:       float64(0),
			expected:  float64(3.14),
		},
		{
			primitive: value(typed.NewNumberValue(int64(42))),
			dst:       int32(0),
			expected:  int32(42),
		},
		{
			primitive: value(typed.NewNumberValue(int64(3))),
			dst:       int16(0),
			expected:  int16(3),
		},
		{
			primitive: value(typed.NewStringValue("hello")),
			dst:       string(""),
			expected:  "hello",
		},
	}

	for _, tc := range testCases {
		if err := tc.primitive.As(&tc.dst); err != nil || !reflect.DeepEqual(tc.dst, tc.expected) {
			t.Fail()
		}
	}
}
