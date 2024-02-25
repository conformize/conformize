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

func TestCollectionHasNoDuplicateElementsPredicate(t *testing.T) {
	tests := []struct {
		name    string
		value   typed.Valuable
		want    bool
		wantErr bool
	}{
		{
			name: "returns true when list has unique elements",
			value: typed.NewListValue(
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
			name: "returns false when list has duplicate elements ",
			value: typed.NewListValue(
				[]typed.Valuable{
					tests.PrimVal(1, &typed.NumberTyped{}),
					tests.PrimVal(2, &typed.NumberTyped{}),
					tests.PrimVal(1, &typed.NumberTyped{}),
				}, &typed.NumberTyped{},
			),
			want:    false,
			wantErr: false,
		},
		{
			name: "returns true when list is empty",
			value: typed.NewListValue(
				[]typed.Valuable{}, &typed.NumberTyped{},
			),
			want:    true,
			wantErr: false,
		},
		{
			name: "returns true when list has elements as unique nested lists",
			value: typed.NewListValue(
				[]typed.Valuable{
					typed.NewListValue(
						[]typed.Valuable{
							tests.PrimVal(1, &typed.NumberTyped{}),
							tests.PrimVal(2, &typed.NumberTyped{}),
						}, &typed.NumberTyped{},
					),
					typed.NewListValue(
						[]typed.Valuable{
							tests.PrimVal(3, &typed.NumberTyped{}),
							tests.PrimVal(4, &typed.NumberTyped{}),
						}, &typed.NumberTyped{},
					),
				}, &typed.ListTyped{ElementsType: &typed.ListTyped{ElementsType: &typed.NumberTyped{}}},
			),
			want:    true,
			wantErr: false,
		},
		{
			name: "returns false when list has elements as duplicate nested lists",
			value: typed.NewListValue(
				[]typed.Valuable{
					typed.NewListValue(
						[]typed.Valuable{
							tests.PrimVal(1, &typed.NumberTyped{}),
							tests.PrimVal(2, &typed.NumberTyped{}),
						}, &typed.NumberTyped{},
					),
					typed.NewListValue(
						[]typed.Valuable{
							tests.PrimVal(1, &typed.NumberTyped{}),
							tests.PrimVal(2, &typed.NumberTyped{}),
						}, &typed.NumberTyped{},
					),
				}, &typed.ListTyped{ElementsType: &typed.ListTyped{ElementsType: &typed.NumberTyped{}}},
			),
			want:    false,
			wantErr: false,
		},
	}

	hasNoDupPrd := &collection.HasNoDuplicateElementsPredicate{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := hasNoDupPrd.Test(tt.value, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("collection.HasNoDuplicateElementsPredicate.Test() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("collection.HasNoDuplicateElementsPredicate.Test() = %v, want %v", got, tt.want)
			}
		})
	}
}
