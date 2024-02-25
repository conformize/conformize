// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package evaluation

import (
	"github.com/conformize/conformize/predicates/predicatefactory"
	"github.com/conformize/conformize/predicates/rule"
)

type RuleEvaluation struct{}

func (c *RuleEvaluation) Evaluate(rule *rule.Rule) (bool, error) {
	predicate, err := predicatefactory.Instance().Build(rule.Predicate)
	if err == nil {
		return predicate.Test(rule.Value, &rule.Arguments)
	}
	return false, err
}

func NewRuleEvaluation() *RuleEvaluation {
	return &RuleEvaluation{}
}
