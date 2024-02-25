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

func TestVariantToStruct(t *testing.T) {
	type TestStructA struct {
		Path      string `cnfrmz:"path"`
		Sensitive bool   `cnfrmz:"sensitive"`
	}

	type TestStructB struct {
		Value     any  `cnfrmz:"value"`
		Sensitive bool `cnfrmz:"sensitive"`
	}

	testCases := []struct {
		variantValue *typed.VariantValue
		value        any
		expected     any
	}{
		{
			variantValue: &typed.VariantValue{
				Value: typed.NewObjectValue(
					map[string]typed.Valuable{
						"path":      value(typed.NewStringValue("$test.path")),
						"sensitive": value(typed.NewBooleanValue(false)),
					},
					map[string]typed.Typeable{
						"path":      &typed.StringTyped{},
						"sensitive": &typed.BooleanTyped{},
					},
				),
				VariantsTypes: []typed.Typeable{
					&typed.ObjectTyped{
						FieldsTypes: map[string]typed.Typeable{
							"path":      &typed.StringTyped{},
							"sensitive": &typed.BooleanTyped{},
						},
					},
					&typed.ObjectTyped{
						FieldsTypes: map[string]typed.Typeable{
							"value":     &typed.GenericTyped{},
							"sensitive": &typed.BooleanTyped{},
						},
					},
				},
			},
			value: &TestStructA{},
			expected: &TestStructA{
				Path:      "$test.path",
				Sensitive: false,
			},
		},
		{
			variantValue: &typed.VariantValue{
				Value: typed.NewObjectValue(
					map[string]typed.Valuable{
						"value":     value(typed.NewStringValue("$test.value")),
						"sensitive": value(typed.NewBooleanValue(false)),
					},
					map[string]typed.Typeable{
						"value":     &typed.StringTyped{},
						"sensitive": &typed.BooleanTyped{},
					},
				),
				VariantsTypes: []typed.Typeable{
					&typed.ObjectTyped{
						FieldsTypes: map[string]typed.Typeable{
							"path":      &typed.StringTyped{},
							"sensitive": &typed.BooleanTyped{},
						},
					},
					&typed.ObjectTyped{
						FieldsTypes: map[string]typed.Typeable{
							"value":     &typed.GenericTyped{},
							"sensitive": &typed.BooleanTyped{},
						},
					},
				},
			},
			value: &TestStructB{},
			expected: &TestStructB{
				Value:     "$test.value",
				Sensitive: false,
			},
		},
	}

	for _, testCase := range testCases {
		if err := testCase.variantValue.As(testCase.value); err != nil || !reflect.DeepEqual(testCase.value, testCase.expected) {
			t.Fail()
		}
	}
}
