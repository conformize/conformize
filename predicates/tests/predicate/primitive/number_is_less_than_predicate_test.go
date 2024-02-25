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

func TestNumberIsLessThanPredicate(t *testing.T) {
	type args struct {
		value typed.Valuable
		args  typed.Valuable
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "returns false when value is equal",
			args: args{
				value: tests.PrimVal(42, &typed.NumberTyped{}),
				args:  tests.PrimVal(42, &typed.NumberTyped{}),
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "returns false when value is greater than argument",
			args: args{
				value: tests.PrimVal(42, &typed.NumberTyped{}),
				args:  tests.PrimVal(3, &typed.NumberTyped{}),
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "returns true when value is less than argument",
			args: args{
				value: tests.PrimVal(3, &typed.NumberTyped{}),
				args:  tests.PrimVal(5, &typed.NumberTyped{}),
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "returns false and an error when value is not a number",
			args: args{
				value: tests.PrimVal("test", &typed.StringTyped{}),
				args:  tests.PrimVal(3, &typed.NumberTyped{}),
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "returns false and an error when argument is not a number",
			args: args{
				value: tests.PrimVal(42, &typed.NumberTyped{}),
				args:  tests.PrimVal("test", &typed.StringTyped{}),
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "returns false and an error when no argument is provided",
			args: args{
				value: tests.PrimVal(42, &typed.NumberTyped{}),
			},
			want:    false,
			wantErr: true,
		},
	}

	numLtPrd := &primitive.NumberIsLessThanPredicate[float64]{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := numLtPrd.Test(tt.args.value, tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("NumberIsLessThanPredicate.Test() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NumberIsLessThanPredicate.Test() = %v, want %v", got, tt.want)
			}
		})
	}
}
