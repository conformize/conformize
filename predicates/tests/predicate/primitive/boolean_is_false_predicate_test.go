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

func TestBooleanIsFalsePredicate(t *testing.T) {
	type args struct {
		value typed.Valuable
	}

	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "returns true when value is false",
			args: args{
				value: tests.PrimVal(false, &typed.BooleanTyped{}),
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "returns false when value is true",
			args: args{
				value: tests.PrimVal(true, &typed.BooleanTyped{}),
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "returns false and error when value is not a boolean",
			args: args{
				value: tests.PrimVal(42, &typed.NumberTyped{}),
			},
			want:    false,
			wantErr: true,
		},
	}

	boolIsFalsePrd := &primitive.BooleanIsFalsePredicate{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("%s", tt.name)
			got, err := boolIsFalsePrd.Test(tt.args.value, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("BooleanIsFalsePredicate.Test() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("BooleanIsFalsePredicate.Test() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBooleanIsFalsePredicate_Arguments(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		{
			name: "returns 0 as number of expected arguments",
			want: 0,
		},
	}

	boolIsFalsePrd := &primitive.BooleanIsFalsePredicate{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := boolIsFalsePrd.Arguments(); got != tt.want {
				t.Errorf("BooleanIsFalsePredicate.Arguments() = %v, want %v", got, tt.want)
			}
		})
	}
}
