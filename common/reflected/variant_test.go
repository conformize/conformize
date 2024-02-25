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
		Path      string
		Sensitive bool
	}

	type TestStructB struct {
		Value     interface{}
		Sensitive bool
	}

	testCases := []struct {
		variantValue *typed.VariantValue
		value        interface{}
		expected     interface{}
	}{
		{
			variantValue: &typed.VariantValue{
				Value: typed.NewObjectValue(
					map[string]typed.Valuable{
						"Path":      value(typed.NewStringValue("$test.path")),
						"Sensitive": value(typed.NewBooleanValue(false)),
					},
					map[string]typed.Typeable{
						"Path":      &typed.StringTyped{},
						"Sensitive": &typed.BooleanTyped{},
					},
				),
				VariantsTypes: []typed.Typeable{
					&typed.ObjectTyped{
						FieldsTypes: map[string]typed.Typeable{
							"Path":      &typed.StringTyped{},
							"Sensitive": &typed.BooleanTyped{},
						},
					},
					&typed.ObjectTyped{
						FieldsTypes: map[string]typed.Typeable{
							"Value":     &typed.GenericTyped{},
							"Sensitive": &typed.BooleanTyped{},
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
	}

	for _, testCase := range testCases {
		if err := testCase.variantValue.As(testCase.value); err != nil || !reflect.DeepEqual(testCase.value, testCase.expected) {
			t.Fail()
		}
	}
}
