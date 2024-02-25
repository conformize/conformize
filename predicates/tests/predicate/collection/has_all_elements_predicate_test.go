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
	"github.com/conformize/conformize/predicates/predicate/collection"
	"github.com/conformize/conformize/predicates/tests"
)

func TestCollectionHasAllElementsPredicate(t *testing.T) {
	tests := []struct {
		name    string
		value   typed.Valuable
		args    typed.Valuable
		want    bool
		wantErr bool
	}{
		{
			name: "returns true when all elements are found in list",
			value: typed.NewListValue(
				[]typed.Valuable{
					tests.PrimVal(42, &typed.NumberTyped{}),
					tests.PrimVal(420, &typed.NumberTyped{}),
					tests.PrimVal(4200, &typed.NumberTyped{}),
					tests.PrimVal(80, &typed.NumberTyped{}),
					tests.PrimVal(8000, &typed.NumberTyped{}),
					tests.PrimVal(8080, &typed.NumberTyped{}),
				}, &typed.NumberTyped{},
			),
			args: &typed.ListValue{
				Elements: []typed.Valuable{
					tests.PrimVal(42, &typed.NumberTyped{}),
					tests.PrimVal(80, &typed.NumberTyped{}),
					tests.PrimVal(8080, &typed.NumberTyped{}),
				},
				ElementsType: &typed.NumberTyped{},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "returns false when not all elements are found in list",
			value: typed.NewListValue(
				[]typed.Valuable{
					tests.PrimVal(42, &typed.NumberTyped{}),
					tests.PrimVal(4200, &typed.NumberTyped{}),
					tests.PrimVal(80, &typed.NumberTyped{}),
					tests.PrimVal(8000, &typed.NumberTyped{}),
					tests.PrimVal(8080, &typed.NumberTyped{}),
				}, &typed.NumberTyped{},
			),
			args: &typed.ListValue{
				Elements: []typed.Valuable{
					tests.PrimVal(42, &typed.NumberTyped{}),
					tests.PrimVal(420, &typed.NumberTyped{}),
					tests.PrimVal(8000, &typed.NumberTyped{}),
				},
				ElementsType: &typed.NumberTyped{},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "returns true when all nested list elements are found",
			value: typed.NewListValue(
				[]typed.Valuable{
					typed.NewListValue(
						[]typed.Valuable{
							typed.NewListValue(
								[]typed.Valuable{
									tests.PrimVal(3.14, &typed.NumberTyped{}),
									tests.PrimVal(42, &typed.NumberTyped{}),
								}, &typed.NumberTyped{},
							),
						}, &typed.ListTyped{ElementsType: &typed.NumberTyped{}},
					),
					typed.NewListValue(
						[]typed.Valuable{
							typed.NewListValue(
								[]typed.Valuable{
									tests.PrimVal(443, &typed.NumberTyped{}),
									tests.PrimVal(8080, &typed.NumberTyped{}),
								}, &typed.NumberTyped{},
							),
						}, &typed.ListTyped{ElementsType: &typed.NumberTyped{}},
					),
					typed.NewListValue(
						[]typed.Valuable{
							typed.NewListValue(
								[]typed.Valuable{
									tests.PrimVal(420, &typed.NumberTyped{}),
									tests.PrimVal(8000, &typed.NumberTyped{}),
								}, &typed.NumberTyped{},
							),
						}, &typed.ListTyped{ElementsType: &typed.NumberTyped{}},
					),
				}, &typed.ListTyped{ElementsType: &typed.ListTyped{ElementsType: &typed.ListTyped{ElementsType: &typed.NumberTyped{}}}},
			),
			args: &typed.ListValue{
				Elements: []typed.Valuable{
					typed.NewListValue(
						[]typed.Valuable{
							typed.NewListValue(
								[]typed.Valuable{
									tests.PrimVal(3.14, &typed.NumberTyped{}),
									tests.PrimVal(42, &typed.NumberTyped{}),
								}, &typed.NumberTyped{},
							),
						}, &typed.ListTyped{ElementsType: &typed.ListTyped{ElementsType: &typed.NumberTyped{}}},
					),
					typed.NewListValue(
						[]typed.Valuable{
							typed.NewListValue(
								[]typed.Valuable{
									tests.PrimVal(443, &typed.NumberTyped{}),
									tests.PrimVal(8080, &typed.NumberTyped{}),
								}, &typed.NumberTyped{},
							),
						}, &typed.ListTyped{ElementsType: &typed.ListTyped{ElementsType: &typed.NumberTyped{}}},
					),
				},
				ElementsType: &typed.ListTyped{ElementsType: &typed.ListTyped{ElementsType: &typed.NumberTyped{}}},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "returns false when not all nested list elements are found list",
			value: typed.NewListValue(
				[]typed.Valuable{
					typed.NewListValue(
						[]typed.Valuable{
							typed.NewListValue(
								[]typed.Valuable{
									tests.PrimVal("hello", &typed.StringTyped{}),
									tests.PrimVal("world", &typed.StringTyped{}),
								}, &typed.StringTyped{},
							),
						}, &typed.ListTyped{ElementsType: &typed.StringTyped{}},
					),
					typed.NewListValue(
						[]typed.Valuable{
							typed.NewListValue(
								[]typed.Valuable{
									tests.PrimVal("test", &typed.StringTyped{}),
									tests.PrimVal("another test", &typed.StringTyped{}),
								}, &typed.StringTyped{},
							),
						}, &typed.ListTyped{ElementsType: &typed.StringTyped{}},
					),
					typed.NewListValue(
						[]typed.Valuable{
							typed.NewListValue(
								[]typed.Valuable{
									tests.PrimVal("yet another test string", &typed.StringTyped{}),
									tests.PrimVal("last test string", &typed.StringTyped{}),
								}, &typed.StringTyped{},
							),
						}, &typed.ListTyped{ElementsType: &typed.StringTyped{}},
					),
				}, &typed.ListTyped{ElementsType: &typed.ListTyped{ElementsType: &typed.ListTyped{ElementsType: &typed.StringTyped{}}}},
			),
			args: &typed.ListValue{
				Elements: []typed.Valuable{
					typed.NewListValue(
						[]typed.Valuable{
							typed.NewListValue(
								[]typed.Valuable{
									tests.PrimVal("test", &typed.StringTyped{}),
									tests.PrimVal("another test", &typed.StringTyped{}),
								}, &typed.StringTyped{},
							),
						}, &typed.ListTyped{ElementsType: &typed.ListTyped{ElementsType: &typed.StringTyped{}}},
					),
					typed.NewListValue(
						[]typed.Valuable{
							typed.NewListValue(
								[]typed.Valuable{
									tests.PrimVal("first not matched", &typed.StringTyped{}),
									tests.PrimVal("second not matched", &typed.StringTyped{}),
								}, &typed.StringTyped{},
							),
						}, &typed.ListTyped{ElementsType: &typed.ListTyped{ElementsType: &typed.StringTyped{}}},
					),
				},
				ElementsType: &typed.ListTyped{ElementsType: &typed.ListTyped{ElementsType: &typed.StringTyped{}}},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "returns false when list is empty",
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

	hasAllPrd := &collection.HasAllElementsPredicate{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := hasAllPrd.Test(tt.value, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("collection.HasAllElementsPredicate.Test() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("collection.HasAllElementsPredicate.Test() = %v, want %v", got, tt.want)
			}
		})
	}
}
