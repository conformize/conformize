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
	"github.com/conformize/conformize/predicates/predicatefactory"
	"github.com/conformize/conformize/predicates/tests"
)

func TestStringHasLengthPredicate(t *testing.T) {
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
			name: "returns true when value has length equal to argument and condition is equal",
			args: args{
				value: tests.PrimVal("hello", &typed.StringTyped{}),
				args: &typed.TupleValue{
					Elements: []typed.Valuable{
						tests.PrimVal("equal", &typed.StringTyped{}),
						tests.PrimVal(5, &typed.NumberTyped{}),
					},
					ElementsTypes: []typed.Typeable{&typed.StringTyped{}, &typed.NumberTyped{}},
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "returns false when value has greater length than argument and condition is equal",
			args: args{
				value: tests.PrimVal("hello", &typed.StringTyped{}),
				args: &typed.TupleValue{
					Elements: []typed.Valuable{
						tests.PrimVal("equal", &typed.StringTyped{}),
						tests.PrimVal(3, &typed.NumberTyped{}),
					},
					ElementsTypes: []typed.Typeable{&typed.StringTyped{}, &typed.NumberTyped{}},
				},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "returns false when value has length less than argument and condition is equal",
			args: args{
				value: tests.PrimVal("test", &typed.StringTyped{}),
				args: &typed.TupleValue{
					Elements: []typed.Valuable{
						tests.PrimVal("equal", &typed.StringTyped{}),
						tests.PrimVal(5, &typed.NumberTyped{}),
					},
					ElementsTypes: []typed.Typeable{&typed.StringTyped{}, &typed.NumberTyped{}},
				},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "returns false and an error when value is not a string and condition is equal",
			args: args{
				value: tests.PrimVal(42, &typed.NumberTyped{}),
				args: &typed.TupleValue{
					Elements: []typed.Valuable{
						tests.PrimVal("equal", &typed.StringTyped{}),
						tests.PrimVal(3, &typed.NumberTyped{}),
					},
					ElementsTypes: []typed.Typeable{&typed.StringTyped{}, &typed.NumberTyped{}},
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "returns false and an error when argument is not a number and condition is equal",
			args: args{
				value: tests.PrimVal("test", &typed.StringTyped{}),
				args: &typed.TupleValue{
					Elements: []typed.Valuable{
						tests.PrimVal("equal", &typed.StringTyped{}),
						tests.PrimVal("test", &typed.StringTyped{}),
					},
					ElementsTypes: []typed.Typeable{&typed.StringTyped{}, &typed.NumberTyped{}},
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "returns false and an error when no argument is provided and condition is equal",
			args: args{
				value: tests.PrimVal(42, &typed.StringTyped{}),
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "returns false when value has equal length and condition is lessThan",
			args: args{
				value: tests.PrimVal("hello", &typed.StringTyped{}),
				args: &typed.TupleValue{
					Elements: []typed.Valuable{
						tests.PrimVal("lessThan", &typed.StringTyped{}),
						tests.PrimVal(5, &typed.NumberTyped{}),
					},
					ElementsTypes: []typed.Typeable{&typed.StringTyped{}, &typed.NumberTyped{}},
				},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "returns false when value has greater length and condition is lessThan",
			args: args{
				value: tests.PrimVal("test", &typed.StringTyped{}),
				args: &typed.TupleValue{
					Elements: []typed.Valuable{
						tests.PrimVal("lessThan", &typed.StringTyped{}),
						tests.PrimVal(3, &typed.NumberTyped{}),
					},
					ElementsTypes: []typed.Typeable{&typed.StringTyped{}, &typed.NumberTyped{}},
				},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "returns true when value has length less than expected and condition is lessThan",
			args: args{
				value: tests.PrimVal("test", &typed.StringTyped{}),
				args: &typed.TupleValue{
					Elements: []typed.Valuable{
						tests.PrimVal("lessThan", &typed.StringTyped{}),
						tests.PrimVal(5, &typed.NumberTyped{}),
					},
					ElementsTypes: []typed.Typeable{&typed.StringTyped{}, &typed.NumberTyped{}},
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "returns false and an error when value is not a string and condition is lessThan",
			args: args{
				value: tests.PrimVal(42, &typed.NumberTyped{}),
				args: &typed.TupleValue{
					Elements: []typed.Valuable{
						tests.PrimVal("lessThan", &typed.StringTyped{}),
						tests.PrimVal(3, &typed.NumberTyped{}),
					},
					ElementsTypes: []typed.Typeable{&typed.StringTyped{}, &typed.StringTyped{}},
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "returns false and an error when argument is not a number and condition is lessThan",
			args: args{
				value: tests.PrimVal("test", &typed.StringTyped{}),
				args: &typed.TupleValue{
					Elements: []typed.Valuable{
						tests.PrimVal("lessThan", &typed.StringTyped{}),
						tests.PrimVal("test", &typed.StringTyped{}),
					},
					ElementsTypes: []typed.Typeable{&typed.StringTyped{}, &typed.StringTyped{}},
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "returns false and an error when no argument is provided and condition is lessThan",
			args: args{
				value: tests.PrimVal(42, &typed.StringTyped{}),
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "returns false when value has equal length to argument and condition is greaterThan",
			args: args{
				value: tests.PrimVal("hello", &typed.StringTyped{}),
				args: &typed.TupleValue{
					Elements: []typed.Valuable{
						tests.PrimVal("greaterThan", &typed.StringTyped{}),
						tests.PrimVal(5, &typed.NumberTyped{}),
					},
					ElementsTypes: []typed.Typeable{&typed.StringTyped{}, &typed.NumberTyped{}},
				},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "returns less when value has length less than argument and condition is greaterThan",
			args: args{
				value: tests.PrimVal("test", &typed.StringTyped{}),
				args: &typed.TupleValue{
					Elements: []typed.Valuable{
						tests.PrimVal("greaterThan", &typed.StringTyped{}),
						tests.PrimVal(5, &typed.NumberTyped{}),
					},
					ElementsTypes: []typed.Typeable{&typed.StringTyped{}, &typed.NumberTyped{}},
				},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "returns true when value has length greater than argument and condition is greaterThan",
			args: args{
				value: tests.PrimVal("test", &typed.StringTyped{}),
				args: &typed.TupleValue{
					Elements: []typed.Valuable{
						tests.PrimVal("greaterThan", &typed.StringTyped{}),
						tests.PrimVal(3, &typed.NumberTyped{}),
					},
					ElementsTypes: []typed.Typeable{&typed.StringTyped{}, &typed.NumberTyped{}},
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "returns false and an error when value is not a string and condition is greaterThan",
			args: args{
				value: tests.PrimVal(42, &typed.NumberTyped{}),
				args: &typed.TupleValue{
					Elements: []typed.Valuable{
						tests.PrimVal("greaterThan", &typed.StringTyped{}),
						tests.PrimVal(3, &typed.NumberTyped{}),
					},
					ElementsTypes: []typed.Typeable{&typed.StringTyped{}, &typed.NumberTyped{}},
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "returns false and an error when argument is not a number and condition is greaterThan",
			args: args{
				value: tests.PrimVal("test", &typed.StringTyped{}),
				args: &typed.TupleValue{
					Elements: []typed.Valuable{
						tests.PrimVal("greaterThan", &typed.StringTyped{}),
						tests.PrimVal("test", &typed.StringTyped{}),
					},
					ElementsTypes: []typed.Typeable{&typed.StringTyped{}, &typed.NumberTyped{}},
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "returns false and an error when no argument is provided and condition is greaterThan",
			args: args{
				value: tests.PrimVal("test", &typed.StringTyped{}),
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "returns true when value has equal length to argument and condition is lessThanOrEqual",
			args: args{
				value: tests.PrimVal("hello", &typed.StringTyped{}),
				args: &typed.TupleValue{
					Elements: []typed.Valuable{
						tests.PrimVal("lessThanOrEqual", &typed.StringTyped{}),
						tests.PrimVal(5, &typed.NumberTyped{}),
					},
					ElementsTypes: []typed.Typeable{&typed.StringTyped{}, &typed.NumberTyped{}},
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "returns true when value has length less than argument and condition is lessThanOrEqual",
			args: args{
				value: tests.PrimVal("test", &typed.StringTyped{}),
				args: &typed.TupleValue{
					Elements: []typed.Valuable{
						tests.PrimVal("lessThanOrEqual", &typed.StringTyped{}),
						tests.PrimVal(5, &typed.NumberTyped{}),
					},
					ElementsTypes: []typed.Typeable{&typed.StringTyped{}, &typed.NumberTyped{}},
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "returns false when value has length greater than argument and condition is lessThanOrEqual",
			args: args{
				value: tests.PrimVal("test", &typed.StringTyped{}),
				args: &typed.TupleValue{
					Elements: []typed.Valuable{
						tests.PrimVal("lessThanOrEqual", &typed.StringTyped{}),
						tests.PrimVal(3, &typed.NumberTyped{}),
					},
					ElementsTypes: []typed.Typeable{&typed.StringTyped{}, &typed.NumberTyped{}},
				},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "returns false and an error when value is not a string and condition is lessThanOrEqual",
			args: args{
				value: tests.PrimVal(42, &typed.NumberTyped{}),
				args: &typed.TupleValue{
					Elements: []typed.Valuable{
						tests.PrimVal("lessThanOrEqual", &typed.StringTyped{}),
						tests.PrimVal(3, &typed.NumberTyped{}),
					},
					ElementsTypes: []typed.Typeable{&typed.StringTyped{}, &typed.NumberTyped{}},
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "returns false and an error when argument is not a number and condition is lessThanOrEqual",
			args: args{
				value: tests.PrimVal("test", &typed.StringTyped{}),
				args: &typed.TupleValue{
					Elements:      []typed.Valuable{tests.PrimVal("test", &typed.StringTyped{})},
					ElementsTypes: []typed.Typeable{&typed.StringTyped{}},
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "returns false and an error when no argument is provided and condition is lessThanOrEqual",
			args: args{
				value: tests.PrimVal("test", &typed.StringTyped{}),
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "returns true when value has equal length to argument and condition is greaterThanOrEqual",
			args: args{
				value: tests.PrimVal("hello", &typed.StringTyped{}),
				args: &typed.TupleValue{
					Elements: []typed.Valuable{
						tests.PrimVal("greaterThanOrEqual", &typed.StringTyped{}),
						tests.PrimVal(5, &typed.NumberTyped{}),
					},
					ElementsTypes: []typed.Typeable{&typed.StringTyped{}, &typed.NumberTyped{}},
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "returns true when value has greater length than argument and condition is greaterThanOrEqual",
			args: args{
				value: tests.PrimVal("test", &typed.StringTyped{}),
				args: &typed.TupleValue{
					Elements: []typed.Valuable{
						tests.PrimVal("greaterThanOrEqual", &typed.StringTyped{}),
						tests.PrimVal(3, &typed.NumberTyped{}),
					},
					ElementsTypes: []typed.Typeable{&typed.StringTyped{}, &typed.NumberTyped{}},
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "returns false when value has length less than argument and condition is greaterThanOrEqual",
			args: args{
				value: tests.PrimVal("test", &typed.StringTyped{}),
				args: &typed.TupleValue{
					Elements: []typed.Valuable{
						tests.PrimVal("greaterThanOrEqual", &typed.StringTyped{}),
						tests.PrimVal(5, &typed.NumberTyped{}),
					},
					ElementsTypes: []typed.Typeable{&typed.StringTyped{}, &typed.NumberTyped{}},
				},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "returns false and an error when value is not a string and condition is greaterThanOrEqual",
			args: args{
				value: tests.PrimVal(42, &typed.NumberTyped{}),
				args: &typed.TupleValue{
					Elements: []typed.Valuable{
						tests.PrimVal("greaterThanOrEqual", &typed.StringTyped{}),
						tests.PrimVal(3, &typed.NumberTyped{}),
					},
					ElementsTypes: []typed.Typeable{&typed.StringTyped{}, &typed.NumberTyped{}},
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "returns false and an error when argument is not a number and condition is greaterThanOrEqual",
			args: args{
				value: tests.PrimVal("test", &typed.StringTyped{}),
				args: &typed.TupleValue{
					Elements: []typed.Valuable{
						tests.PrimVal("greaterThanOrEqual", &typed.StringTyped{}),
						tests.PrimVal("test", &typed.StringTyped{}),
					},
					ElementsTypes: []typed.Typeable{&typed.StringTyped{}, &typed.NumberTyped{}},
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "returns false and an error when no argument is provided and condition is greaterThanOrEqual",
			args: args{
				value: tests.PrimVal("test", &typed.StringTyped{}),
				args: &typed.TupleValue{
					Elements: []typed.Valuable{
						tests.PrimVal("greaterThanOrEqual", &typed.StringTyped{}),
					},
					ElementsTypes: []typed.Typeable{&typed.StringTyped{}},
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "returns true when value has equal length to lower bound argument and condition is withinRange",
			args: args{
				value: tests.PrimVal("hello", &typed.StringTyped{}),
				args: &typed.TupleValue{
					Elements: []typed.Valuable{
						tests.PrimVal("withinRange", &typed.StringTyped{}),
						tests.PrimVal(5, &typed.NumberTyped{}),
						tests.PrimVal(5, &typed.NumberTyped{}),
					},
					ElementsTypes: []typed.Typeable{&typed.StringTyped{}, &typed.NumberTyped{}, &typed.NumberTyped{}},
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "returns true when value has greater length than lower bound argument and condition is withinRange",
			args: args{
				value: tests.PrimVal("test", &typed.StringTyped{}),
				args: &typed.TupleValue{
					Elements: []typed.Valuable{
						tests.PrimVal("withinRange", &typed.StringTyped{}),
						tests.PrimVal(3, &typed.NumberTyped{}),
						tests.PrimVal(5, &typed.NumberTyped{}),
					},
					ElementsTypes: []typed.Typeable{&typed.StringTyped{}, &typed.NumberTyped{}, &typed.NumberTyped{}},
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "returns false when value has length less than lower bound argument and condition is withinRange",
			args: args{
				value: tests.PrimVal("test", &typed.StringTyped{}),
				args: &typed.TupleValue{
					Elements: []typed.Valuable{
						tests.PrimVal("withinRange", &typed.StringTyped{}),
						tests.PrimVal(10, &typed.NumberTyped{}),
						tests.PrimVal(15, &typed.NumberTyped{}),
					},
					ElementsTypes: []typed.Typeable{&typed.StringTyped{}, &typed.NumberTyped{}, &typed.NumberTyped{}},
				},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "returns true when value has equal length to upper bound argument and condition is withinRange",
			args: args{
				value: tests.PrimVal("hello", &typed.StringTyped{}),
				args: &typed.TupleValue{
					Elements: []typed.Valuable{
						tests.PrimVal("withinRange", &typed.StringTyped{}),
						tests.PrimVal(1, &typed.NumberTyped{}),
						tests.PrimVal(5, &typed.NumberTyped{}),
					},
					ElementsTypes: []typed.Typeable{&typed.StringTyped{}, &typed.NumberTyped{}, &typed.NumberTyped{}},
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "returns true when value has lower length than upper bound argument and condition is withinRange",
			args: args{
				value: tests.PrimVal("test", &typed.StringTyped{}),
				args: &typed.TupleValue{
					Elements: []typed.Valuable{
						tests.PrimVal("withinRange", &typed.StringTyped{}),
						tests.PrimVal(3, &typed.NumberTyped{}),
						tests.PrimVal(5, &typed.NumberTyped{}),
					},
					ElementsTypes: []typed.Typeable{&typed.StringTyped{}, &typed.NumberTyped{}, &typed.NumberTyped{}},
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "returns false when value has length greater than upper bound argument and condition is withinRange",
			args: args{
				value: tests.PrimVal("test123", &typed.StringTyped{}),
				args: &typed.TupleValue{
					Elements: []typed.Valuable{
						tests.PrimVal("withinRange", &typed.StringTyped{}),
						tests.PrimVal(1, &typed.NumberTyped{}),
						tests.PrimVal(5, &typed.NumberTyped{}),
					},
					ElementsTypes: []typed.Typeable{&typed.StringTyped{}, &typed.NumberTyped{}, &typed.NumberTyped{}},
				},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "returns false and an error when value is not a string and condition is withinRange",
			args: args{
				value: tests.PrimVal(42, &typed.NumberTyped{}),
				args: &typed.TupleValue{
					Elements: []typed.Valuable{
						tests.PrimVal("withinRange", &typed.StringTyped{}),
						tests.PrimVal(10, &typed.NumberTyped{}),
						tests.PrimVal(15, &typed.NumberTyped{}),
					},
					ElementsTypes: []typed.Typeable{&typed.StringTyped{}, &typed.NumberTyped{}, &typed.NumberTyped{}},
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "returns false and an error when argument is not a numbe and condition is withinRanger",
			args: args{
				value: tests.PrimVal("test", &typed.StringTyped{}),
				args: &typed.TupleValue{
					Elements: []typed.Valuable{
						tests.PrimVal("withinRange", &typed.StringTyped{}),
						tests.PrimVal("test", &typed.StringTyped{}),
						tests.PrimVal("test", &typed.StringTyped{}),
					},
					ElementsTypes: []typed.Typeable{&typed.StringTyped{}, &typed.NumberTyped{}, &typed.NumberTyped{}},
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "returns false and an error when no argument is provided and condition is withinRange",
			args: args{
				value: tests.PrimVal("test", &typed.StringTyped{}),
			},
			want:    false,
			wantErr: true,
		},
	}

	strLenPrd := &primitive.StringHasLengthPredicate{PredicateBuilder: predicatefactory.Instance()}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := strLenPrd.Test(tt.args.value, tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("StringHasLengthPredicate.Test() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("StringHasLengthPredicate.Test() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringHasLengthPredicate_Arguments(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		{
			name: "returns 1 as expected number of arguments",
			want: 1,
		},
	}

	strLenPrd := &primitive.StringHasLengthPredicate{PredicateBuilder: predicatefactory.Instance()}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := strLenPrd.Arguments(); got != tt.want {
				t.Errorf("StringHasLengthPredicate.Arguments() = %v, want %v", got, tt.want)
			}
		})
	}
}
