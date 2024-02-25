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
	"github.com/conformize/conformize/predicates/predicate/primitive"
	"github.com/conformize/conformize/predicates/tests"
)

func TestBooleanIsTruePredicate(t *testing.T) {
	tests := []struct {
		value   typed.Valuable
		name    string
		want    bool
		wantErr bool
	}{
		{
			value:   tests.PrimVal(true, &typed.BooleanTyped{}),
			name:    "returns true when value is true",
			want:    true,
			wantErr: false,
		},
		{
			value:   tests.PrimVal(false, &typed.BooleanTyped{}),
			name:    "returns false when value is false",
			want:    false,
			wantErr: false,
		},
		{
			value:   tests.PrimVal(42, &typed.NumberTyped{}),
			name:    "returns false and error when value is not a boolean",
			want:    false,
			wantErr: true,
		},
	}

	boolIsTruePrd := &primitive.BooleanIsTruePredicate{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := boolIsTruePrd.Test(tt.value, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("BooleanIsTruePredicate.Test() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("BooleanIsTruePredicate.Test() = %v, want %v", got, tt.want)
			}
		})
	}
}
