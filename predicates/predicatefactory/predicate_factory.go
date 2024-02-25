// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package predicatefactory

import (
	"fmt"
	"sync"

	"github.com/conformize/conformize/predicates"
	"github.com/conformize/conformize/predicates/condition"
)

type PredicateFactory struct {
	predicateBuilders *PredicateBuilderFuncsMap
}

func (prdFactory *PredicateFactory) Build(condition condition.ConditionType) (predicates.Predicate, error) {
	predicate, ok := (*prdFactory.predicateBuilders)[condition]
	if ok {
		return predicate(prdFactory), nil
	}
	return nil, fmt.Errorf("predicate for condition %s not found", condition.String())
}

func newPredicateFactory() predicates.PredicateBuilder {
	return &PredicateFactory{predicateBuilders}
}

var instance predicates.PredicateBuilder
var once sync.Once

func Instance() predicates.PredicateBuilder {
	once.Do(func() {
		instance = newPredicateFactory()
	})
	return instance
}
