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

func TestStringMatchesExpressionPredicate(t *testing.T) {
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
			name: "returns true when value matches expression",
			args: args{
				value: tests.PrimVal("192.168.1.1", &typed.StringTyped{}),
				args:  tests.PrimVal("^(\\d{1,3}\\.){3}\\d{1,3}$", &typed.StringTyped{}),
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "returns false when value doesn't match expression",
			args: args{
				value: tests.PrimVal("192.168.1.1", &typed.StringTyped{}),
				args:  tests.PrimVal("^(\\d{1,3}\\.){3}\\d{1,3}$", &typed.StringTyped{}),
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "returns false and error when value is not a string",
			args: args{
				value: tests.PrimVal(42, &typed.NumberTyped{}),
				args:  tests.PrimVal("^(\\d{1,3}\\.){3}\\d{1,3}$", &typed.StringTyped{}),
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "returns false and error when argument is not a string",
			args: args{
				value: tests.PrimVal("192.168.1.1", &typed.StringTyped{}),
				args:  tests.PrimVal(42, &typed.NumberTyped{}),
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "returns false and error when no argument is provided",
			args: args{
				value: tests.PrimVal("192.168.1.1", &typed.StringTyped{}),
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "returns false and error when expression is not valid",
			args: args{
				value: tests.PrimVal("hello", &typed.StringTyped{}),
				args:  tests.PrimVal("[unclosed", &typed.StringTyped{}),
			},
			want:    false,
			wantErr: true,
		},
	}

	strMatchExprPrd := primitive.StringMatchesExpressionPredicate{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := strMatchExprPrd.Test(tt.args.value, tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("StringMatchesExpressionPredicate.Test() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("StringMatchesExpressionPredicate.Test() = %v, want %v", got, tt.want)
			}
		})
	}
}
