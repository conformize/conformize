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
	"github.com/conformize/conformize/predicates/predicate/equality"
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
			name: "returns ValueIsEqualPredicate for condition equal and string value",
			args: args{
				value:     tests.PrimVal("test", &typed.StringTyped{}),
				condition: condition.EQ,
			},
			want: &equality.ValueIsEqualPredicate{PredicateBuilder: Instance()},
		},
		{
			name: "returns StringMatchesExpressionPredicate for condition matchesExpression and string value",
			args: args{
				value:     tests.PrimVal("test", &typed.StringTyped{}),
				condition: condition.MATCHES,
			},
			want: &primitive.StringMatchesExpressionPredicate{},
		},
		{
			name: "returns ValueIsEqualPredicate for condition equal and number value",
			args: args{
				value:     tests.PrimVal(42, &typed.NumberTyped{}),
				condition: condition.EQ,
			},
			want: &equality.ValueIsEqualPredicate{PredicateBuilder: Instance()},
		},
		{
			name: "returns ValueIsEqualPredicate for condition greaterThan and number value",
			args: args{
				value:     tests.PrimVal(42, &typed.NumberTyped{}),
				condition: condition.GT,
			},
			want: &primitive.NumberIsGreaterThanPredicate[float64]{},
		},
		{
			name: "returns NumberIsLessThanPredicate for condition lessThan and number value",
			args: args{
				value:     tests.PrimVal(42, &typed.NumberTyped{}),
				condition: condition.LT,
			},
			want: &primitive.NumberIsLessThanPredicate[float64]{},
		},
		{
			name: "returns NumberIsGreaterThanOrEqualPredicate for condition greaterThanOrEqual and number value",
			args: args{
				value:     tests.PrimVal(42, &typed.NumberTyped{}),
				condition: condition.GTE,
			},
			want: &primitive.NumberIsGreaterThanOrEqualPredicate[float64]{},
		},
		{
			name: "returns NumberIsLessThanOrEqualPredicate for condition lessThanOrEqual and number value",
			args: args{
				value:     tests.PrimVal(42, &typed.NumberTyped{}),
				condition: condition.LTE,
			},
			want: &primitive.NumberIsLessThanOrEqualPredicate[float64]{},
		},
		{
			name: "returns NumberIsWithinRangePredicate for condition withinRange and number value",
			args: args{
				value:     tests.PrimVal(42, &typed.NumberTyped{}),
				condition: condition.RANGE,
			},
			want: &primitive.NumberIsWithinRangePredicate[float64]{},
		},
		{
			name: "returns BooleanIsTruePredicate for condition isTrue and boolean value",
			args: args{
				value:     tests.PrimVal(true, &typed.BooleanTyped{}),
				condition: condition.TRUE,
			},
			want: &primitive.BooleanIsTruePredicate{},
		},
		{
			name: "returns BooleanIsFalsePredicate for condition isFalse and boolean value",
			args: args{
				value:     tests.PrimVal(true, &typed.BooleanTyped{}),
				condition: condition.FALSE,
			},
			want: &primitive.BooleanIsFalsePredicate{},
		},
		{
			name: "returns BooleanIsEqualPredicate for condition equal and boolean value",
			args: args{
				value:     tests.PrimVal(true, &typed.BooleanTyped{}),
				condition: condition.EQ,
			},
			want: &equality.ValueIsEqualPredicate{PredicateBuilder: Instance()},
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
