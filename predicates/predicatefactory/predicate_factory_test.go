// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package predicatefactory

import (
	"reflect"
	"testing"

	"github.com/conformize/conformize/common/typed"
	"github.com/conformize/conformize/predicates"
	"github.com/conformize/conformize/predicates/condition"
	"github.com/conformize/conformize/predicates/predicate/primitive"
	"github.com/conformize/conformize/predicates/tests"
)

func TestPredicateFactory_Build(t *testing.T) {
	type args struct {
		value     typed.Valuable
		condition condition.ConditionType
	}
	tests := []struct {
		name    string
		args    args
		want    predicates.Predicate
		wantErr bool
	}{
		{
			name: "returns StringIsEqualPredicate for condition EQUAL and string value",
			args: args{
				value:     tests.PrimVal("test", &typed.StringTyped{}),
				condition: condition.EQUAL,
			},
			want: &primitive.StringIsEqualPredicate{},
		},
		{
			name: "returns StringMatchesExpressionPredicate for condition MATCHES_EXPRESISON and string value",
			args: args{
				value:     tests.PrimVal("test", &typed.StringTyped{}),
				condition: condition.MATCHES_EXPRESSION,
			},
			want: &primitive.StringMatchesExpressionPredicate{},
		},
		{
			name: "returns NumberIsEqualPredicate for condition EQUAL and number value",
			args: args{
				value:     tests.PrimVal(42, &typed.NumberTyped{}),
				condition: condition.EQUAL,
			},
			want: &primitive.NumberIsEqualPredicate[float64]{},
		},
		{
			name: "returns NumberIsGreaterThanPredicate for condition GREATER_THAN and number value",
			args: args{
				value:     tests.PrimVal(42, &typed.NumberTyped{}),
				condition: condition.GREATER_THAN,
			},
			want: &primitive.NumberIsGreaterThanPredicate[float64]{},
		},
		{
			name: "returns NumberIsLessThanPredicate for condition LESS_THAN and number value",
			args: args{
				value:     tests.PrimVal(42, &typed.NumberTyped{}),
				condition: condition.LESS_THAN,
			},
			want: &primitive.NumberIsLessThanPredicate[float64]{},
		},
		{
			name: "returns NumberIsGreaterThanOrEqualPredicate for condition GREATER_THAN_OR_EQUAL and number value",
			args: args{
				value:     tests.PrimVal(42, &typed.NumberTyped{}),
				condition: condition.GREATER_THAN_OR_EQUAL,
			},
			want: &primitive.NumberIsGreaterThanOrEqualPredicate[float64]{},
		},
		{
			name: "returns NumberIsLessThanOrEqualPredicate for condition LESS_THAN_OR_EQUAL and number value",
			args: args{
				value:     tests.PrimVal(42, &typed.NumberTyped{}),
				condition: condition.LESS_THAN_OR_EQUAL,
			},
			want: &primitive.NumberIsLessThanOrEqualPredicate[float64]{},
		},
		{
			name: "returns NumberIsWithinRangePredicate for condition WITHIN_RANGE and number value",
			args: args{
				value:     tests.PrimVal(42, &typed.NumberTyped{}),
				condition: condition.WITHIN_RANGE,
			},
			want: &primitive.NumberIsWithinRangePredicate[float64]{},
		},
		{
			name: "returns BooleanIsTruePredicate for condition IS_TRUE and boolean value",
			args: args{
				value:     tests.PrimVal(true, &typed.BooleanTyped{}),
				condition: condition.IS_TRUE,
			},
			want: &primitive.BooleanIsTruePredicate{},
		},
		{
			name: "returns BooleanIsFalsePredicate for condition IS_FALSE and boolean value",
			args: args{
				value:     tests.PrimVal(true, &typed.BooleanTyped{}),
				condition: condition.IS_FALSE,
			},
			want: &primitive.BooleanIsFalsePredicate{},
		},
		{
			name: "returns BooleanIsEqualPredicate for condition EQUAL and boolean value",
			args: args{
				value:     tests.PrimVal(true, &typed.BooleanTyped{}),
				condition: condition.EQUAL,
			},
			want: &primitive.BooleanIsEqualPredicate{},
		},
		{
			name: "returns nil predicate and error when conditon is not matched for value",
			args: args{
				value:     tests.PrimVal(true, &typed.BooleanTyped{}),
				condition: condition.HAS_LENGTH,
			},
			want:    nil,
			wantErr: true,
		},
	}
	prdFactory := Instance()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := prdFactory.Build(tt.args.condition)
			if (err != nil) != tt.wantErr {
				t.Errorf("PredicateFactory.Build() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PredicateFactory.Build() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInstance(t *testing.T) {
	tests := []struct {
		name string
		want predicates.PredicateBuilder
	}{
		{
			name: "Test PredicateFactory is initialized",
			want: &PredicateFactory{predicateBuilders: predicateBuilders},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Instance(); got == nil || !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Instance() = %v, want %v", got, tt.want)
			}
		})
	}
}
