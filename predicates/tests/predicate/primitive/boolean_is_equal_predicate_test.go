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

func TestBooleanIsEqualPredicate(t *testing.T) {
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
			name: "returns true when boolean values are equal",
			args: args{
				value: tests.PrimVal(true, &typed.BooleanTyped{}),
				args:  tests.PrimVal(true, &typed.BooleanTyped{}),
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "returns false when boolean values are not equal",
			args: args{
				value: tests.PrimVal(true, &typed.BooleanTyped{}),
				args:  tests.PrimVal(false, &typed.BooleanTyped{}),
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "returns false and error when value is not a boolean",
			args: args{
				value: tests.PrimVal(42, &typed.NumberTyped{}),
				args:  tests.PrimVal(true, &typed.BooleanTyped{}),
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "returns false and error when argument is not a boolean",
			args: args{
				value: tests.PrimVal(true, &typed.BooleanTyped{}),
				args: &typed.TupleValue{
					Elements: []typed.Valuable{
						tests.PrimVal("True", &typed.StringTyped{}),
					},
					ElementsTypes: []typed.Typeable{&typed.BooleanTyped{}},
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "returns false and error whehn argument is not provided",
			args: args{
				value: tests.PrimVal(42, &typed.NumberTyped{}),
				args:  nil,
			},
			want:    false,
			wantErr: true,
		},
	}

	boolEqPrd := &primitive.BooleanIsEqualPredicate{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := boolEqPrd.Test(tt.args.value, tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("BooleanIsEqualPredicate.Test() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("BooleanIsEqualPredicate.Test() = %v, want %v", got, tt.want)
			}
		})
	}
}
