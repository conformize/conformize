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
	"github.com/conformize/conformize/predicates/tests"
)

func TestListContainsAllElementsPredicate(t *testing.T) {
	tests := []struct {
		name    string
		value   typed.Valuable
		args    *typed.TupleValue
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
			args: &typed.TupleValue{
				Elements: []typed.Valuable{
					&typed.ListValue{
						Elements: []typed.Valuable{
							tests.PrimVal(42, &typed.NumberTyped{}),
							tests.PrimVal(80, &typed.NumberTyped{}),
							tests.PrimVal(8080, &typed.NumberTyped{}),
						},
						ElementsType: &typed.NumberTyped{},
					},
				},
				ElementsTypes: []typed.Typeable{&typed.ListTyped{ElementsType: &typed.NumberTyped{}}},
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
			args: &typed.TupleValue{
				Elements: []typed.Valuable{
					&typed.ListValue{
						Elements: []typed.Valuable{
							tests.PrimVal(42, &typed.NumberTyped{}),
							tests.PrimVal(420, &typed.NumberTyped{}),
							tests.PrimVal(8000, &typed.NumberTyped{}),
						},
						ElementsType: &typed.NumberTyped{},
					},
				},
				ElementsTypes: []typed.Typeable{&typed.ListTyped{ElementsType: &typed.NumberTyped{}}},
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
			args: &typed.TupleValue{
				Elements: []typed.Valuable{
					&typed.ListValue{
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
				},
				ElementsTypes: []typed.Typeable{
					&typed.ListTyped{ElementsType: &typed.ListTyped{ElementsType: &typed.ListTyped{ElementsType: &typed.NumberTyped{}}}},
				},
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
			args: &typed.TupleValue{
				Elements: []typed.Valuable{
					&typed.ListValue{
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
				},
				ElementsTypes: []typed.Typeable{&typed.ListTyped{ElementsType: &typed.ListTyped{ElementsType: &typed.StringTyped{}}}},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "returns false when list is empty",
			value: typed.NewListValue(
				[]typed.Valuable{}, &typed.BooleanTyped{},
			),
			args: &typed.TupleValue{
				Elements: []typed.Valuable{
					&typed.ListValue{
						Elements:     []typed.Valuable{tests.PrimVal(true, &typed.BooleanTyped{})},
						ElementsType: &typed.BooleanTyped{},
					},
				},
				ElementsTypes: []typed.Typeable{&typed.ListTyped{ElementsType: &typed.BooleanTyped{}}},
			},
			want:    false,
			wantErr: false,
		},
	}

	listContainsAllPrd := &list.ListContainsAllElementsPredicate{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := listContainsAllPrd.Test(tt.value, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListContainsAllElementsPredicate.Test() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("TListContainsAllElementsPredicate.Test() = %v, want %v", got, tt.want)
			}
		})
	}
}