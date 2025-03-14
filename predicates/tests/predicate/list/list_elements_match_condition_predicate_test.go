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
	"github.com/conformize/conformize/predicates"
	"github.com/conformize/conformize/predicates/predicate/list"
	"github.com/conformize/conformize/predicates/predicatefactory"
	"github.com/conformize/conformize/predicates/tests"
)

func TestListElementsMatchConditionPredicate(t *testing.T) {
	tests := []struct {
		name      string
		value     typed.Valuable
		predicate predicates.Predicate
		args      *typed.TupleValue
		want      bool
		wantErr   bool
	}{
		{
			name: "test match all elements witch condition greaterThan returns true when condition is met",
			predicate: &list.ListElementsMatchConditionPredicate{
				PredicateBuilder: predicatefactory.Instance(),
			},
			value: typed.NewListValue(
				[]typed.Valuable{
					tests.PrimVal(6, &typed.NumberTyped{}),
					tests.PrimVal(15, &typed.NumberTyped{}),
					tests.PrimVal(16, &typed.NumberTyped{}),
				}, &typed.NumberTyped{},
			),
			args: &typed.TupleValue{
				Elements:      []typed.Valuable{tests.PrimVal("greaterThan", &typed.StringTyped{}), tests.PrimVal(5, &typed.NumberTyped{})},
				ElementsTypes: []typed.Typeable{&typed.StringTyped{}, &typed.NumberTyped{}},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "test match all elements witch condition lessThan returns false when condition is not met",
			predicate: &list.ListElementsMatchConditionPredicate{
				PredicateBuilder: predicatefactory.Instance(),
			},
			value: typed.NewListValue(
				[]typed.Valuable{
					tests.PrimVal(101, &typed.NumberTyped{}),
					tests.PrimVal(102, &typed.NumberTyped{}),
					tests.PrimVal(103, &typed.NumberTyped{}),
				}, &typed.NumberTyped{},
			),
			args: &typed.TupleValue{
				Elements:      []typed.Valuable{tests.PrimVal("lessThan", &typed.StringTyped{}), tests.PrimVal(42, &typed.NumberTyped{})},
				ElementsTypes: []typed.Typeable{&typed.StringTyped{}, &typed.NumberTyped{}},
			},
			want:    false,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.predicate.Test(tt.value, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListElementsMatchConditionPredicate.Test() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ListElementsMatchConditionPredicate.Test() = %v, want %v", got, tt.want)
			}
		})
	}
}
