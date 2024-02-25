// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package evaluation

import (
	"fmt"

	"github.com/conformize/conformize/common"
	"github.com/conformize/conformize/common/typed"
	"github.com/conformize/conformize/predicates/rule"
)

type RuleEvaluation struct{}

func (c *RuleEvaluation) Evaluate(rule *rule.Rule) (bool, error) {
	if val, ok := rule.Value.(typed.Valuable); ok {
		return rule.Predicate.Test(val)
	}

	if fnVal, ok := rule.Value.(*common.IterFnNodeValue); ok {
		return fnVal.Fn(fnVal.Iter, rule.Predicate)
	}

	return false, fmt.Errorf("rule evaluation: unsupported value type %T", rule.Value)
}

func NewRuleEvaluation() *RuleEvaluation {
	return &RuleEvaluation{}
}
