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

func TestCollectionIsSubsetPredicate(t *testing.T) {
	tests := []struct {
		name    string
		value   typed.Valuable
		args    typed.Valuable
		want    bool
		wantErr bool
	}{
		{
			name: "returns true when list is a subset",
			value: typed.NewListValue(
				[]typed.Valuable{
					tests.PrimVal(1, &typed.NumberTyped{}),
					tests.PrimVal(2, &typed.NumberTyped{}),
				}, &typed.NumberTyped{},
			),
			args: typed.NewListValue(
				[]typed.Valuable{
					tests.PrimVal(1, &typed.NumberTyped{}),
					tests.PrimVal(2, &typed.NumberTyped{}),
					tests.PrimVal(3, &typed.NumberTyped{}),
				}, &typed.NumberTyped{},
			),
			want:    true,
			wantErr: false,
		},
		{
			name: "returns false when list is not a subset",
			value: typed.NewListValue(
				[]typed.Valuable{
					tests.PrimVal(4, &typed.NumberTyped{}),
					tests.PrimVal(5, &typed.NumberTyped{}),
				}, &typed.NumberTyped{},
			),
			args: typed.NewListValue(
				[]typed.Valuable{
					tests.PrimVal(1, &typed.NumberTyped{}),
					tests.PrimVal(2, &typed.NumberTyped{}),
					tests.PrimVal(3, &typed.NumberTyped{}),
				}, &typed.NumberTyped{},
			),
			want:    false,
			wantErr: false,
		},
		{
			name: "returns true when a list with nested lists is a subset",
			value: typed.NewListValue(
				[]typed.Valuable{
					typed.NewListValue(
						[]typed.Valuable{
							tests.PrimVal(1, &typed.NumberTyped{}),
							tests.PrimVal(3, &typed.NumberTyped{}),
							tests.PrimVal(5, &typed.NumberTyped{}),
						}, &typed.NumberTyped{},
					),
				}, &typed.ListTyped{ElementsType: &typed.ListTyped{ElementsType: &typed.NumberTyped{}}},
			),
			args: typed.NewListValue(
				[]typed.Valuable{
					typed.NewListValue(
						[]typed.Valuable{
							tests.PrimVal(1, &typed.NumberTyped{}),
							tests.PrimVal(3, &typed.NumberTyped{}),
							tests.PrimVal(5, &typed.NumberTyped{}),
						}, &typed.NumberTyped{},
					),
					typed.NewListValue(
						[]typed.Valuable{
							tests.PrimVal(2, &typed.NumberTyped{}),
							tests.PrimVal(4, &typed.NumberTyped{}),
							tests.PrimVal(6, &typed.NumberTyped{}),
						}, &typed.NumberTyped{},
					),
				}, &typed.ListTyped{ElementsType: &typed.ListTyped{ElementsType: &typed.NumberTyped{}}},
			),
			want:    true,
			wantErr: false,
		},
		{
			name: "returns false when a list with nested lists is not a subset",
			value: typed.NewListValue(
				[]typed.Valuable{
					typed.NewListValue(
						[]typed.Valuable{
							tests.PrimVal(4, &typed.NumberTyped{}),
							tests.PrimVal(5, &typed.NumberTyped{}),
						}, &typed.NumberTyped{},
					),
				}, &typed.ListTyped{ElementsType: &typed.ListTyped{ElementsType: &typed.NumberTyped{}}},
			),
			args: typed.NewListValue(
				[]typed.Valuable{
					typed.NewListValue(
						[]typed.Valuable{
							tests.PrimVal(1, &typed.NumberTyped{}),
							tests.PrimVal(2, &typed.NumberTyped{}),
							tests.PrimVal(3, &typed.NumberTyped{}),
						}, &typed.NumberTyped{},
					),
				}, &typed.ListTyped{ElementsType: &typed.ListTyped{ElementsType: &typed.NumberTyped{}}},
			),
			want:    false,
			wantErr: false,
		},
		{
			name: "returns true when list is empty",
			value: typed.NewListValue(
				[]typed.Valuable{}, &typed.NumberTyped{},
			),
			args: typed.NewListValue(
				[]typed.Valuable{
					tests.PrimVal(1, &typed.NumberTyped{}),
					tests.PrimVal(2, &typed.NumberTyped{}),
					tests.PrimVal(3, &typed.NumberTyped{}),
				}, &typed.NumberTyped{},
			),
			want:    true,
			wantErr: false,
		},
		{
			name: "returns false when superset list is empty",
			value: typed.NewListValue(
				[]typed.Valuable{
					tests.PrimVal(1, &typed.NumberTyped{}),
				}, &typed.NumberTyped{},
			),
			args:    typed.NewListValue([]typed.Valuable{}, &typed.NumberTyped{}),
			want:    false,
			wantErr: false,
		},
		{
			name: "returns true when both subset and superset lists are empty",
			value: typed.NewListValue(
				[]typed.Valuable{}, &typed.NumberTyped{},
			),
			args:    typed.NewListValue([]typed.Valuable{}, &typed.NumberTyped{}),
			want:    true,
			wantErr: false,
		},
		{
			name: "returns false when subset and superset lists have different types",
			value: typed.NewListValue(
				[]typed.Valuable{
					tests.PrimVal(1, &typed.NumberTyped{}),
				}, &typed.NumberTyped{},
			),
			args: &typed.TupleValue{
				Elements: []typed.Valuable{typed.NewListValue([]typed.Valuable{
					tests.PrimVal("string", &typed.StringTyped{}),
				}, &typed.NumberTyped{})},
				ElementsTypes: []typed.Typeable{&typed.ListTyped{ElementsType: &typed.StringTyped{}}},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "returns true when list with repeating elements is a subset",
			value: typed.NewListValue(
				[]typed.Valuable{
					tests.PrimVal(1, &typed.NumberTyped{}),
					tests.PrimVal(1, &typed.NumberTyped{}),
				}, &typed.NumberTyped{},
			),
			args: typed.NewListValue(
				[]typed.Valuable{
					tests.PrimVal(1, &typed.NumberTyped{}),
					tests.PrimVal(1, &typed.NumberTyped{}),
					tests.PrimVal(2, &typed.NumberTyped{}),
				}, &typed.NumberTyped{},
			),
			want:    true,
			wantErr: false,
		},
		{
			name: "returns true when repeating elements in a list have more occurences than in superset",
			value: typed.NewListValue(
				[]typed.Valuable{
					tests.PrimVal(1, &typed.NumberTyped{}),
					tests.PrimVal(1, &typed.NumberTyped{}),
					tests.PrimVal(1970, &typed.NumberTyped{}),
				}, &typed.NumberTyped{},
			),
			args: typed.NewListValue(
				[]typed.Valuable{
					tests.PrimVal(1, &typed.NumberTyped{}),
					tests.PrimVal(1, &typed.NumberTyped{}),
					tests.PrimVal(1, &typed.NumberTyped{}),
					tests.PrimVal(1970, &typed.NumberTyped{}),
				}, &typed.NumberTyped{},
			),
			want:    true,
			wantErr: false,
		},
		{
			name: "returns false when repeating elements in a list have less occurences than in superset",
			value: typed.NewListValue(
				[]typed.Valuable{
					tests.PrimVal(1, &typed.NumberTyped{}),
					tests.PrimVal(1, &typed.NumberTyped{}),
					tests.PrimVal(1, &typed.NumberTyped{}),
					tests.PrimVal(4, &typed.NumberTyped{}),
				}, &typed.NumberTyped{},
			),
			args: typed.NewListValue(
				[]typed.Valuable{
					tests.PrimVal(1, &typed.NumberTyped{}),
					tests.PrimVal(1, &typed.NumberTyped{}),
					tests.PrimVal(4, &typed.NumberTyped{}),
				}, &typed.NumberTyped{},
			),
			want:    false,
			wantErr: false,
		},
		{
			name: "returns true when deeply nested list is a subset",
			value: typed.NewListValue(
				[]typed.Valuable{
					typed.NewListValue(
						[]typed.Valuable{
							typed.NewListValue(
								[]typed.Valuable{
									tests.PrimVal(1, &typed.NumberTyped{}),
									tests.PrimVal(1, &typed.NumberTyped{}),
									tests.PrimVal(1970, &typed.NumberTyped{}),
								}, &typed.NumberTyped{},
							),
							typed.NewListValue(
								[]typed.Valuable{
									tests.PrimVal(1, &typed.NumberTyped{}),
									tests.PrimVal(1, &typed.NumberTyped{}),
									tests.PrimVal(1983, &typed.NumberTyped{}),
								}, &typed.NumberTyped{},
							),
						}, &typed.ListTyped{ElementsType: &typed.NumberTyped{}},
					),
				}, &typed.ListTyped{ElementsType: &typed.ListTyped{ElementsType: &typed.NumberTyped{}}},
			),
			args: typed.NewListValue(
				[]typed.Valuable{
					typed.NewListValue(
						[]typed.Valuable{
							typed.NewListValue(
								[]typed.Valuable{
									tests.PrimVal(1, &typed.NumberTyped{}),
									tests.PrimVal(1, &typed.NumberTyped{}),
									tests.PrimVal(1970, &typed.NumberTyped{}),
								}, &typed.NumberTyped{},
							),
							typed.NewListValue(
								[]typed.Valuable{
									tests.PrimVal(1, &typed.NumberTyped{}),
									tests.PrimVal(1, &typed.NumberTyped{}),
									tests.PrimVal(1983, &typed.NumberTyped{}),
								}, &typed.NumberTyped{},
							),
						}, &typed.ListTyped{ElementsType: &typed.NumberTyped{}},
					),
					typed.NewListValue(
						[]typed.Valuable{
							typed.NewListValue(
								[]typed.Valuable{
									tests.PrimVal(0, &typed.NumberTyped{}),
									tests.PrimVal(1, &typed.NumberTyped{}),
									tests.PrimVal(1, &typed.NumberTyped{}),
								}, &typed.NumberTyped{},
							),
							typed.NewListValue(
								[]typed.Valuable{
									tests.PrimVal(2, &typed.NumberTyped{}),
									tests.PrimVal(3, &typed.NumberTyped{}),
									tests.PrimVal(5, &typed.NumberTyped{}),
								}, &typed.NumberTyped{},
							),
						}, &typed.ListTyped{ElementsType: &typed.NumberTyped{}},
					),
				}, &typed.ListTyped{ElementsType: &typed.ListTyped{ElementsType: &typed.ListTyped{ElementsType: &typed.NumberTyped{}}}},
			),
			want:    true,
			wantErr: false,
		},
	}

	isSubsetPrd := &collection.IsSubsetPredicate{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := isSubsetPrd.Test(tt.value, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("collection.IsSubsetPredicate.Test() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("collection.IsSubsetPredicate.Test() = %v, want %v", got, tt.want)
			}
		})
	}
}
