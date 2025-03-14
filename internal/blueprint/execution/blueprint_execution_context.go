// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package execution

import (
	"github.com/conformize/conformize/internal/blueprint/elements"
	"github.com/conformize/conformize/predicates/ruleset"
)

type BlueprintExecutionContext struct {
	rulesMeta []*elements.RuleMeta
	ruleSet   *ruleset.RuleSet
}

func NewBlueprintExecutionContext() *BlueprintExecutionContext {
	return &BlueprintExecutionContext{
		ruleSet: ruleset.NewEmptyRuleSet(),
	}
}
