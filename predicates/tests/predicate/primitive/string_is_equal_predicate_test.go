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

func TestStringIsEqualPredicate(t *testing.T) {
	type args struct {
		value typed.Valuable
		args  *typed.TupleValue
	}

	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "returns true when value and argument are equal",
			args: args{
				value: tests.PrimVal("test", &typed.StringTyped{}),
				args: &typed.TupleValue{
					Elements:      []typed.Valuable{tests.PrimVal("test", &typed.StringTyped{})},
					ElementsTypes: []typed.Typeable{&typed.StringTyped{}},
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "returns false when value and argument are not equal",
			args: args{
				value: tests.PrimVal("test", &typed.StringTyped{}),
				args: &typed.TupleValue{
					Elements:      []typed.Valuable{tests.PrimVal("hello", &typed.StringTyped{})},
					ElementsTypes: []typed.Typeable{&typed.StringTyped{}},
				},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "returns false and error when value is not a string",
			args: args{
				value: tests.PrimVal(42, &typed.NumberTyped{}),
				args: &typed.TupleValue{
					Elements:      []typed.Valuable{tests.PrimVal("hello", &typed.StringTyped{})},
					ElementsTypes: []typed.Typeable{&typed.StringTyped{}},
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "returns false and error when argument is not a string",
			args: args{
				value: tests.PrimVal("hello", &typed.StringTyped{}),
				args: &typed.TupleValue{
					Elements:      []typed.Valuable{tests.PrimVal(42, &typed.NumberTyped{})},
					ElementsTypes: []typed.Typeable{&typed.NumberTyped{}},
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "returns false and error when no argument is provided",
			args: args{
				value: tests.PrimVal("hello", &typed.StringTyped{}),
			},
			want:    false,
			wantErr: true,
		},
	}

	strLenPrd := &primitive.StringIsEqualPredicate{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := strLenPrd.Test(tt.args.value, tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("StringIsEqualPredicate.Test() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("StringIsEqualPredicate.Test() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringIsEqualPredicate_Arguments(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		{
			name: "returns 1 as expected number of arguments",
			want: 1,
		},
	}

	strLenPrd := &primitive.StringIsEqualPredicate{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := strLenPrd.Arguments(); got != tt.want {
				t.Errorf("StringIsEqualPredicate.Arguments() = %v, want %v", got, tt.want)
			}
		})
	}
}
