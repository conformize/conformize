// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package ruleset

import (
	"github.com/conformize/conformize/predicates/rule"
)

type RuleSet struct {
	Rules []*rule.Rule
}

func (r *RuleSet) AddRule(rule *rule.Rule) {
	r.Rules = append(r.Rules, rule)
}

func NewRuleSet(rules []*rule.Rule) *RuleSet {
	return &RuleSet{Rules: rules}
}

func NewEmptyRuleSet() *RuleSet {
	return &RuleSet{Rules: make([]*rule.Rule, 0)}
}
