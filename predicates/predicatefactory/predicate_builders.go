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
	"github.com/conformize/conformize/predicates/predicate/date"
	"github.com/conformize/conformize/predicates/predicate/equality"
	"github.com/conformize/conformize/predicates/predicate/length"
	"github.com/conformize/conformize/predicates/predicate/list"
	"github.com/conformize/conformize/predicates/predicate/primitive"
)

type PredicateBuilderFunc func(prdBuilder predicates.PredicateBuilder) predicates.Predicate
type PredicateBuilderFuncsMap map[condition.ConditionType]PredicateBuilderFunc

var predicateBuilders = &PredicateBuilderFuncsMap{
	condition.EQUAL: func(prdBuilder predicates.PredicateBuilder) predicates.Predicate {
		return &equality.ValueIsEqualPredicate{PredicateBuilder: prdBuilder}
	},
	condition.NOT_EQUAL: func(prdBuilder predicates.PredicateBuilder) predicates.Predicate {
		return &equality.ValueIsNotEqualPredicate{}
	},
	condition.HAS_LENGTH: func(prdBuilder predicates.PredicateBuilder) predicates.Predicate {
		return &length.ValueHasLengthPredicate{}
	},
	condition.MATCHES_EXPRESSION: func(prdBuilder predicates.PredicateBuilder) predicates.Predicate {
		return &primitive.StringMatchesExpressionPredicate{}
	},
	condition.GREATER_THAN: func(prdBuilder predicates.PredicateBuilder) predicates.Predicate {
		return &primitive.NumberIsGreaterThanPredicate[float64]{}
	},
	condition.LESS_THAN: func(prdBuilder predicates.PredicateBuilder) predicates.Predicate {
		return &primitive.NumberIsLessThanPredicate[float64]{}
	},
	condition.GREATER_THAN_OR_EQUAL: func(prdBuilder predicates.PredicateBuilder) predicates.Predicate {
		return &primitive.NumberIsGreaterThanOrEqualPredicate[float64]{}
	},
	condition.LESS_THAN_OR_EQUAL: func(prdBuilder predicates.PredicateBuilder) predicates.Predicate {
		return &primitive.NumberIsLessThanOrEqualPredicate[float64]{}
	},
	condition.WITHIN_RANGE: func(prdBuilder predicates.PredicateBuilder) predicates.Predicate {
		return &primitive.NumberIsWithinRangePredicate[float64]{}
	},
	condition.IS_TRUE: func(prdBuilder predicates.PredicateBuilder) predicates.Predicate {
		return &primitive.BooleanIsTruePredicate{}
	},
	condition.IS_FALSE: func(prdBuilder predicates.PredicateBuilder) predicates.Predicate {
		return &primitive.BooleanIsFalsePredicate{}
	},
	condition.IS_EMPTY: func(prdBuilder predicates.PredicateBuilder) predicates.Predicate {
		return &list.ListIsEmptyPredicate{}
	},
	condition.MATCH: func(prdBuilder predicates.PredicateBuilder) predicates.Predicate {
		return &list.ListElementsMatchConditionPredicate{PredicateBuilder: prdBuilder}
	},
	condition.ANY_MATCH: func(prdBuilder predicates.PredicateBuilder) predicates.Predicate {
		return &list.ListAnyElementMatcheshConditionPredicate{PredicateBuilder: prdBuilder}
	},
	condition.NO_MATCH: func(prdBuilder predicates.PredicateBuilder) predicates.Predicate {
		return &list.ListNoElementMatcheConditionPredicate{PredicateBuilder: prdBuilder}
	},
	condition.NO_DUPLICATES: func(prdBuilder predicates.PredicateBuilder) predicates.Predicate {
		return &list.ListHasNoDuplicateElementsPredicate{}
	},
	condition.CONTAINS: func(prdBuilder predicates.PredicateBuilder) predicates.Predicate {
		return &list.ListContainsElementPredicate{PredicateBuilder: prdBuilder}
	},
	condition.CONTAINS_ANY_OF: func(prdBuilder predicates.PredicateBuilder) predicates.Predicate {
		return &list.ListContainsAnyOfElementsPredicate{PredicateBuilder: prdBuilder}
	},
	condition.CONTAINS_ALL_OF: func(prdBuilder predicates.PredicateBuilder) predicates.Predicate {
		return &list.ListContainsAllElementsPredicate{PredicateBuilder: prdBuilder}
	},
	condition.IS_SUBSET: func(prdBuilder predicates.PredicateBuilder) predicates.Predicate {
		return &list.ListIsSubsetPredicate{}
	},
	condition.SAME_DATE: func(prdBuilder predicates.PredicateBuilder) predicates.Predicate {
		return &date.DateIsPredicate{}
	},
	condition.NOT_SAME_DATE: func(prdBuilder predicates.PredicateBuilder) predicates.Predicate {
		return &date.DateIsNotPredicate{}
	},
	condition.VALID_DATE: func(prdBuilder predicates.PredicateBuilder) predicates.Predicate {
		return &date.DateIsValidPredicate{}
	},
	condition.DATE_UP_TO: func(prdBuilder predicates.PredicateBuilder) predicates.Predicate {
		return &date.DateUpToPredicate{}
	},
	condition.DATE_FROM: func(prdBuilder predicates.PredicateBuilder) predicates.Predicate {
		return &date.DateFromPredicate{}
	},
	condition.DATE_BEFORE: func(prdBuilder predicates.PredicateBuilder) predicates.Predicate {
		return &date.DateIsBeforePredicate{}
	},
	condition.DATE_AFTER: func(prdBuilder predicates.PredicateBuilder) predicates.Predicate {
		return &date.DateIsAfterPredicate{}
	},
	condition.DATE_WITHIN_INTERVAL: func(prdBuilder predicates.PredicateBuilder) predicates.Predicate {
		return &date.DateIsWithinIntervalPredicate{}
	},
	condition.DATE_IN_FUTURE: func(prdBuilder predicates.PredicateBuilder) predicates.Predicate {
		return &date.DateIsInFuture{}
	},
}
