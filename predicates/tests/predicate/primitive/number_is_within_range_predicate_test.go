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

func TestNumberIsWithinRangePredicate(t *testing.T) {
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
			name: "returns true when value is within range",
			args: args{
				value: tests.PrimVal(70, &typed.NumberTyped{}),
				args: &typed.ListValue{
					Elements:     []typed.Valuable{tests.PrimVal(42, &typed.NumberTyped{}), tests.PrimVal(80, &typed.NumberTyped{})},
					ElementsType: &typed.NumberTyped{},
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "returns true when value is equal to first argument",
			args: args{
				value: tests.PrimVal(42, &typed.NumberTyped{}),
				args: &typed.ListValue{
					Elements:     []typed.Valuable{tests.PrimVal(42, &typed.NumberTyped{}), tests.PrimVal(80, &typed.NumberTyped{})},
					ElementsType: &typed.NumberTyped{},
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "returns true when value is equal to second argument",
			args: args{
				value: tests.PrimVal(80, &typed.NumberTyped{}),
				args: &typed.ListValue{
					Elements:     []typed.Valuable{tests.PrimVal(42, &typed.NumberTyped{}), tests.PrimVal(80, &typed.NumberTyped{})},
					ElementsType: &typed.NumberTyped{},
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "returns false when value is less than first argument",
			args: args{
				value: tests.PrimVal(3, &typed.NumberTyped{}),
				args: &typed.ListValue{
					Elements:     []typed.Valuable{tests.PrimVal(42, &typed.NumberTyped{}), tests.PrimVal(80, &typed.NumberTyped{})},
					ElementsType: &typed.NumberTyped{},
				},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "returns false when value is greater than second argument",
			args: args{
				value: tests.PrimVal(8080, &typed.NumberTyped{}),
				args: &typed.ListValue{
					Elements:     []typed.Valuable{tests.PrimVal(42, &typed.NumberTyped{}), tests.PrimVal(80, &typed.NumberTyped{})},
					ElementsType: &typed.NumberTyped{},
				},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "returns false and an error when value is not a number",
			args: args{
				value: tests.PrimVal("test", &typed.StringTyped{}),
				args: &typed.ListValue{
					Elements:     []typed.Valuable{tests.PrimVal(42, &typed.NumberTyped{}), tests.PrimVal(80, &typed.NumberTyped{})},
					ElementsType: &typed.NumberTyped{},
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "returns false and an error when first argument is not a number",
			args: args{
				value: tests.PrimVal(42, &typed.NumberTyped{}),
				args: &typed.ListValue{
					Elements:     []typed.Valuable{tests.PrimVal("test", &typed.StringTyped{}), tests.PrimVal(80, &typed.NumberTyped{})},
					ElementsType: &typed.NumberTyped{},
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "returns false and an error when second argument is not a number",
			args: args{
				value: tests.PrimVal(42, &typed.NumberTyped{}),
				args: &typed.ListValue{
					Elements:     []typed.Valuable{tests.PrimVal(80, &typed.NumberTyped{}), tests.PrimVal("test", &typed.StringTyped{})},
					ElementsType: &typed.NumberTyped{},
				},
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
		{
			name: "returns false and an error when insufficient number of arguments are provided",
			args: args{
				value: tests.PrimVal(42, &typed.NumberTyped{}),
				args: &typed.ListValue{
					Elements:     []typed.Valuable{tests.PrimVal(80, &typed.NumberTyped{})},
					ElementsType: &typed.NumberTyped{},
				},
			},
			want:    false,
			wantErr: true,
		},
	}

	numWithinRangePrd := &primitive.NumberIsWithinRangePredicate[float64]{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := numWithinRangePrd.Test(tt.args.value, tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("NumberIsWithinRangePredicate.Test() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NumberIsWithinRangePredicate.Test() = %v, want %v", got, tt.want)
			}
		})
	}
}
