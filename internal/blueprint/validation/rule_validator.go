// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package validation

import (
	"fmt"

	"github.com/conformize/conformize/internal/blueprint/elements"
	"github.com/conformize/conformize/predicates/condition"
)

type ruleValidator struct {
	ruleArgumentValidator
}

func (v *ruleValidator) Validate(rule *elements.Rule, configSources map[string]elements.ConfigurationSource, refs map[string]string) error {
	if len(rule.Value.Steps()) == 0 {
		return fmt.Errorf("value path not specified")
	}
	valPathSteps := rule.Value.Steps()
	root := valPathSteps[0].String()
	if _, sourceRefFound := configSources[root]; !sourceRefFound {
		if _, aliasRefFound := refs[root]; !aliasRefFound {
			return fmt.Errorf("couldn't resolve root '%s' in path %s", root, rule.Value.String())
		}
	}

	predicate := condition.FromString(rule.Predicate)
	if predicate == condition.UNKNOWN {
		return fmt.Errorf("couldn't resolve predicate '%s'", rule.Predicate)
	}

	if err := v.ruleArgumentValidator.Validate(rule.Arguments, configSources, refs); err != nil {
		return err
	}
	return nil
}
