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

func TestListContainsElementPredicate(t *testing.T) {
	tests := []struct {
		name    string
		value   typed.Valuable
		args    *typed.TupleValue
		want    bool
		wantErr bool
	}{
		{
			name: "returns true when element is found in list",
			value: typed.NewListValue(
				[]typed.Valuable{
					tests.PrimVal("blablabla", &typed.StringTyped{}),
					tests.PrimVal("blabla", &typed.StringTyped{}),
					tests.PrimVal("hello", &typed.StringTyped{}),
				}, &typed.StringTyped{},
			),
			args: &typed.TupleValue{
				Elements:      []typed.Valuable{tests.PrimVal("hello", &typed.StringTyped{})},
				ElementsTypes: []typed.Typeable{&typed.StringTyped{}},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "returns false when element is not found in list",
			value: typed.NewListValue(
				[]typed.Valuable{
					tests.PrimVal("blablabla", &typed.StringTyped{}),
					tests.PrimVal("blabla", &typed.StringTyped{}),
					tests.PrimVal("hello", &typed.StringTyped{}),
				}, &typed.StringTyped{},
			),
			args: &typed.TupleValue{
				Elements:      []typed.Valuable{tests.PrimVal("world", &typed.StringTyped{})},
				ElementsTypes: []typed.Typeable{&typed.StringTyped{}},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "returns true when nested list element is found",
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
				}, &typed.ListTyped{ElementsType: &typed.ListTyped{ElementsType: &typed.ListTyped{ElementsType: &typed.NumberTyped{}}}},
			),
			args: &typed.TupleValue{
				Elements: []typed.Valuable{
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
				},
				ElementsTypes: []typed.Typeable{&typed.ListTyped{ElementsType: &typed.NumberTyped{}}},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "returns false when nested list element is not found",
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
											tests.PrimVal(3200, &typed.NumberTyped{}),
											tests.PrimVal(4200, &typed.NumberTyped{}),
										}, &typed.NumberTyped{},
									),
								}, &typed.ListTyped{ElementsType: &typed.ListTyped{ElementsType: &typed.NumberTyped{}}},
							),
						},
						ElementsType: &typed.ListTyped{ElementsType: &typed.ListTyped{ElementsType: &typed.ListTyped{ElementsType: &typed.NumberTyped{}}}},
					},
				},
				ElementsTypes: []typed.Typeable{&typed.ListTyped{ElementsType: &typed.ListTyped{ElementsType: &typed.ListTyped{ElementsType: &typed.NumberTyped{}}}}},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "returns false when list is empty",
			value: typed.NewListValue(
				[]typed.Valuable{}, &typed.StringTyped{},
			),
			args: &typed.TupleValue{
				Elements:      []typed.Valuable{tests.PrimVal("nothing", &typed.StringTyped{})},
				ElementsTypes: []typed.Typeable{&typed.StringTyped{}},
			},
			want:    false,
			wantErr: false,
		},
	}

	listContainsPrd := &list.ListContainsElementPredicate{PredicateBuilder: predicatefactory.Instance()}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := listContainsPrd.Test(tt.value, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListContainsElementPredicate.Test() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ListContainsElementPredicate.Test() = %v, want %v", got, tt.want)
			}
		})
	}
}
