// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package predicatefactory

import (
	"github.com/conformize/conformize/predicates"
	"github.com/conformize/conformize/predicates/condition"
	"github.com/conformize/conformize/predicates/predicate/collection"
	"github.com/conformize/conformize/predicates/predicate/date"
	"github.com/conformize/conformize/predicates/predicate/equality"
	"github.com/conformize/conformize/predicates/predicate/primitive"
)

type PredicateBuilderFunc func(prdBuilder predicates.PredicateBuilder) predicates.Predicate
type PredicateBuilderFuncsMap map[condition.ConditionType]PredicateBuilderFunc

var predicateBuilders = &PredicateBuilderFuncsMap{
	condition.EQ: func(prdBuilder predicates.PredicateBuilder) predicates.Predicate {
		return &equality.ValueIsEqualPredicate{PredicateBuilder: prdBuilder}
	},
	condition.NOT: func(prdBuilder predicates.PredicateBuilder) predicates.Predicate {
		return &equality.ValueIsNotEqualPredicate{}
	},
	condition.MATCHES: func(prdBuilder predicates.PredicateBuilder) predicates.Predicate {
		return &primitive.StringMatchesExpressionPredicate{}
	},
	condition.GT: func(prdBuilder predicates.PredicateBuilder) predicates.Predicate {
		return &primitive.NumberIsGreaterThanPredicate[float64]{}
	},
	condition.LT: func(prdBuilder predicates.PredicateBuilder) predicates.Predicate {
		return &primitive.NumberIsLessThanPredicate[float64]{}
	},
	condition.GTE: func(prdBuilder predicates.PredicateBuilder) predicates.Predicate {
		return &primitive.NumberIsGreaterThanOrEqualPredicate[float64]{}
	},
	condition.LTE: func(prdBuilder predicates.PredicateBuilder) predicates.Predicate {
		return &primitive.NumberIsLessThanOrEqualPredicate[float64]{}
	},
	condition.RANGE: func(prdBuilder predicates.PredicateBuilder) predicates.Predicate {
		return &primitive.NumberIsWithinRangePredicate[float64]{}
	},
	condition.TRUE: func(prdBuilder predicates.PredicateBuilder) predicates.Predicate {
		return &primitive.BooleanIsTruePredicate{}
	},
	condition.FALSE: func(prdBuilder predicates.PredicateBuilder) predicates.Predicate {
		return &primitive.BooleanIsFalsePredicate{}
	},
	condition.EMPTY: func(prdBuilder predicates.PredicateBuilder) predicates.Predicate {
		return &collection.IsEmptyPredicate{}
	},
	condition.UNIQUE: func(prdBuilder predicates.PredicateBuilder) predicates.Predicate {
		return &collection.HasNoDuplicateElementsPredicate{}
	},
	condition.HAS_ANY: func(prdBuilder predicates.PredicateBuilder) predicates.Predicate {
		return &collection.HasAnyOfElementsPredicate{PredicateBuilder: prdBuilder}
	},
	condition.HAS: func(prdBuilder predicates.PredicateBuilder) predicates.Predicate {
		return &collection.HasAllElementsPredicate{PredicateBuilder: prdBuilder}
	},
	condition.SUBSET_OF: func(prdBuilder predicates.PredicateBuilder) predicates.Predicate {
		return &collection.IsSubsetPredicate{}
	},
	condition.SAME: func(prdBuilder predicates.PredicateBuilder) predicates.Predicate {
		return &date.DateIsPredicate{}
	},
	condition.DIFFERENT: func(prdBuilder predicates.PredicateBuilder) predicates.Predicate {
		return &date.DateIsNotPredicate{}
	},
	condition.VALID: func(prdBuilder predicates.PredicateBuilder) predicates.Predicate {
		return &date.DateIsValidPredicate{}
	},
	condition.UNTIL: func(prdBuilder predicates.PredicateBuilder) predicates.Predicate {
		return &date.DateUpToPredicate{}
	},
	condition.SINCE: func(prdBuilder predicates.PredicateBuilder) predicates.Predicate {
		return &date.DateFromPredicate{}
	},
	condition.BEFORE: func(prdBuilder predicates.PredicateBuilder) predicates.Predicate {
		return &date.DateIsBeforePredicate{}
	},
	condition.AFTER: func(prdBuilder predicates.PredicateBuilder) predicates.Predicate {
		return &date.DateIsAfterPredicate{}
	},
	condition.WITHIN: func(prdBuilder predicates.PredicateBuilder) predicates.Predicate {
		return &date.DateIsWithinIntervalPredicate{}
	},
	condition.FUTURE: func(prdBuilder predicates.PredicateBuilder) predicates.Predicate {
		return &date.DateIsInFuture{}
	},
}
