// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package tests

import (
	"testing"

	"github.com/conformize/conformize/common/typed"
	"github.com/conformize/conformize/predicates/predicate/list"
	"github.com/conformize/conformize/predicates/predicatefactory"
	"github.com/conformize/conformize/predicates/tests"
)

func TestListContainsAnyOfElementsPredicate(t *testing.T) {
	tests := []struct {
		name    string
		value   typed.Valuable
		args    typed.Valuable
		want    bool
		wantErr bool
	}{
		{
			name: "returns true when any of elements are found in list",
			value: typed.NewListValue(
				[]typed.Valuable{
					tests.PrimVal("true", &typed.StringTyped{}),
					tests.PrimVal("42", &typed.StringTyped{}),
					tests.PrimVal("hello", &typed.StringTyped{}),
				}, &typed.StringTyped{},
			),
			args: &typed.ListValue{
				Elements:     []typed.Valuable{tests.PrimVal("hello", &typed.StringTyped{}), tests.PrimVal("world", &typed.StringTyped{})},
				ElementsType: &typed.StringTyped{},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "returns false when noneOf of elements are found in list",
			value: typed.NewListValue(
				[]typed.Valuable{
					tests.PrimVal("true", &typed.StringTyped{}),
					tests.PrimVal("42", &typed.StringTyped{}),
					tests.PrimVal("world", &typed.StringTyped{}),
				}, &typed.StringTyped{},
			),
			args: &typed.ListValue{
				Elements:     []typed.Valuable{tests.PrimVal("test", &typed.StringTyped{}), tests.PrimVal("blablabla", &typed.StringTyped{})},
				ElementsType: &typed.StringTyped{},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "returns true when any of nested list elements are found in list",
			value: typed.NewListValue(
				[]typed.Valuable{
					typed.NewListValue(
						[]typed.Valuable{
							tests.PrimVal(124, &typed.NumberTyped{}),
							tests.PrimVal(42, &typed.NumberTyped{}),
							tests.PrimVal(4300, &typed.NumberTyped{}),
						}, &typed.NumberTyped{},
					),
					typed.NewListValue(
						[]typed.Valuable{
							tests.PrimVal(124, &typed.NumberTyped{}),
							tests.PrimVal(42, &typed.NumberTyped{}),
							tests.PrimVal(4300, &typed.NumberTyped{}),
						}, &typed.BooleanTyped{},
					),
				}, &typed.ListTyped{ElementsType: &typed.NumberTyped{}},
			),
			args: &typed.ListValue{
				Elements: []typed.Valuable{
					typed.NewListValue(
						[]typed.Valuable{
							tests.PrimVal(111, &typed.NumberTyped{}),
							tests.PrimVal(101, &typed.NumberTyped{}),
							tests.PrimVal(10101, &typed.NumberTyped{}),
						}, &typed.NumberTyped{},
					),
					typed.NewListValue(
						[]typed.Valuable{
							tests.PrimVal(124, &typed.NumberTyped{}),
							tests.PrimVal(42, &typed.NumberTyped{}),
							tests.PrimVal(4300, &typed.NumberTyped{}),
						}, &typed.NumberTyped{},
					),
				},
				ElementsType: &typed.ListTyped{ElementsType: &typed.ListTyped{ElementsType: &typed.NumberTyped{}}},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "returns false when noneOf of nested list elements are found in list",
			value: typed.NewListValue(
				[]typed.Valuable{
					typed.NewListValue(
						[]typed.Valuable{
							tests.PrimVal(124, &typed.NumberTyped{}),
							tests.PrimVal(42, &typed.NumberTyped{}),
							tests.PrimVal(4300, &typed.NumberTyped{}),
						}, &typed.NumberTyped{},
					),
					typed.NewListValue(
						[]typed.Valuable{
							tests.PrimVal(124, &typed.NumberTyped{}),
							tests.PrimVal(42, &typed.NumberTyped{}),
							tests.PrimVal(4300, &typed.NumberTyped{}),
						}, &typed.BooleanTyped{},
					),
				}, &typed.ListTyped{ElementsType: &typed.NumberTyped{}},
			),
			args: &typed.ListValue{
				Elements: []typed.Valuable{
					typed.NewListValue(
						[]typed.Valuable{
							tests.PrimVal(11, &typed.NumberTyped{}),
							tests.PrimVal(10, &typed.NumberTyped{}),
							tests.PrimVal(101, &typed.NumberTyped{}),
						}, &typed.NumberTyped{},
					),
					typed.NewListValue(
						[]typed.Valuable{
							tests.PrimVal(24, &typed.NumberTyped{}),
							tests.PrimVal(40, &typed.NumberTyped{}),
							tests.PrimVal(00, &typed.NumberTyped{}),
						}, &typed.NumberTyped{},
					),
				},
				ElementsType: &typed.ListTyped{ElementsType: &typed.ListTyped{ElementsType: &typed.NumberTyped{}}},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "return false when list is empty",
			value: typed.NewListValue(
				[]typed.Valuable{}, &typed.BooleanTyped{},
			),
			args: &typed.ListValue{
				Elements:     []typed.Valuable{tests.PrimVal(true, &typed.BooleanTyped{})},
				ElementsType: &typed.BooleanTyped{},
			},
			want:    false,
			wantErr: false,
		},
	}

	listContainsAnyPrd := &list.ListContainsAnyOfElementsPredicate{PredicateBuilder: predicatefactory.Instance()}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			listContainsAnyPrd.Arguments(tt.args)
			got, err := listContainsAnyPrd.Test(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListContainsAnyOfElementsPredicate.Test() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ListContainsAnyOfElementsPredicate.Test() = %v, want %v", got, tt.want)
			}
		})
	}
}
