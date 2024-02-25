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

func TestMapValue(t *testing.T) {
	testCases := []struct {
		valueType typed.Typeable
		value     any
		expected  *typed.MapValue
	}{
		{
			valueType: &typed.MapTyped{ElementsType: &typed.NumberTyped{}},
			value:     map[string]int64{"1": 2, "3": 4},
			expected: &typed.MapValue{
				ElementsType: &typed.NumberTyped{},
				Elements: map[string]typed.Valuable{
					"1": value(typed.NewNumberValue(2)),
					"3": value(typed.NewNumberValue(4)),
				},
			},
		},
		{
			valueType: &typed.MapTyped{ElementsType: &typed.StringTyped{}},
			value:     map[string]string{"hello": "world"},
			expected: &typed.MapValue{
				ElementsType: &typed.StringTyped{},
				Elements:     map[string]typed.Valuable{"hello": value(typed.NewStringValue("world"))},
			},
		},
		{
			valueType: &typed.MapTyped{ElementsType: &typed.BooleanTyped{}},
			value:     map[string]bool{"true": true, "false": false},
			expected: &typed.MapValue{
				ElementsType: &typed.BooleanTyped{},
				Elements: map[string]typed.Valuable{
					"true":  value(typed.NewBooleanValue(true)),
					"false": value(typed.NewBooleanValue(false)),
				},
			},
		},
		{
			valueType: &typed.MapTyped{ElementsType: &typed.NumberTyped{}},
			value:     map[string]float64{"3.14": 2.71},
			expected: &typed.MapValue{
				ElementsType: &typed.NumberTyped{},
				Elements:     map[string]typed.Valuable{"3.14": value(typed.NewNumberValue(2.71))},
			},
		},
	}

	for _, tc := range testCases {
		if mapVal, err := Map(reflect.ValueOf(tc.value), tc.valueType); err != nil || !reflect.DeepEqual(mapVal, tc.expected) {
			t.Fail()
		}
	}
}

func TestMapValueToMapWithStringValues(t *testing.T) {
	var m map[string]string
	var mapValue typed.Valuable = &typed.MapValue{
		ElementsType: &typed.StringTyped{},
		Elements: map[string]typed.Valuable{
			"hello": value(typed.NewStringValue("world")),
		},
	}
	var expected = map[string]string{
		"hello": "world",
	}

	if err := mapValue.As(&m); err != nil || !reflect.DeepEqual(&m, &expected) {
		t.Fail()
	}
}

func TestMapValueToMapWithIntValues(t *testing.T) {
	var m map[string]int16
	var mapValue typed.Valuable = &typed.MapValue{
		ElementsType: &typed.StringTyped{},
		Elements: map[string]typed.Valuable{
			"nodes": value(typed.NewNumberValue(2)),
		},
	}
	var expected = map[string]int16{
		"nodes": 2,
	}

	if err := mapValue.As(&m); err != nil || !reflect.DeepEqual(&m, &expected) {
		t.Fail()
	}
}

func TestNestedMap(t *testing.T) {
	var m map[string]map[string]int16
	var mapValue typed.Valuable = &typed.MapValue{
		ElementsType: &typed.MapTyped{
			ElementsType: &typed.NumberTyped{},
		},
		Elements: map[string]typed.Valuable{
			"nodes": &typed.MapValue{
				ElementsType: &typed.NumberTyped{},
				Elements: map[string]typed.Valuable{
					"1": value(typed.NewNumberValue(2)),
					"3": value(typed.NewNumberValue(4)),
				},
			},
		},
	}
	var expected = map[string]map[string]int16{
		"nodes": {
			"1": 2,
			"3": 4,
		},
	}

	if err := mapValue.As(&m); err != nil || !reflect.DeepEqual(&m, &expected) {
		t.Fail()
	}
}
