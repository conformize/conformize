// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package rule

import (
	"github.com/conformize/conformize/common/typed"
	"github.com/conformize/conformize/predicates/condition"
)

type RuleBuilder struct {
	predicate condition.ConditionType
	value     typed.Valuable
	arguments typed.TupleValue
}

func (ruleBlder *RuleBuilder) Predicate(predicate condition.ConditionType) *RuleBuilder {
	ruleBlder.predicate = predicate
	return ruleBlder
}

func (ruleBlder *RuleBuilder) Value(value typed.Valuable) *RuleBuilder {
	ruleBlder.value = value
	return ruleBlder
}

func (ruleBlder *RuleBuilder) Argument(argument typed.Valuable) *RuleBuilder {
	ruleBlder.arguments.Elements = append(ruleBlder.arguments.Elements, argument)
	ruleBlder.arguments.ElementsTypes = append(ruleBlder.arguments.ElementsTypes, argument.Type())
	return ruleBlder
}

func (ruleBlder *RuleBuilder) Build() *Rule {
	return &Rule{
		Predicate: ruleBlder.predicate,
		Value:     ruleBlder.value,
		Arguments: ruleBlder.arguments,
	}
}

func NewRuleBuilder() *RuleBuilder {
	return &RuleBuilder{}
}
