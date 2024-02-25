// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package rule

import (
	"github.com/conformize/conformize/common/functions"
	"github.com/conformize/conformize/predicates"
)

type RuleBuilder struct {
	predicate predicates.Predicate
	value     interface{}
	arguments interface{}
}

func (ruleBlder *RuleBuilder) Predicate(predicate predicates.Predicate) *RuleBuilder {
	ruleBlder.predicate = predicate
	return ruleBlder
}

func (ruleBlder *RuleBuilder) Value(value interface{}) *RuleBuilder {
	ruleBlder.value = value
	return ruleBlder
}

func (ruleBlder *RuleBuilder) Arguments(arguments interface{}) *RuleBuilder {
	ruleBlder.arguments = arguments
	return ruleBlder
}

func (ruleBlder *RuleBuilder) Build() (*Rule, error) {
	if argsPrd, ok := ruleBlder.predicate.(predicates.ArgumentsPredicate); ok {
		if argVal, err := functions.ParseRawValue(ruleBlder.arguments); err == nil {
			argsPrd.Arguments(argVal)
		} else {
			return nil, err
		}
	}

	return &Rule{
		Predicate: ruleBlder.predicate,
		Value:     ruleBlder.value,
		Arguments: ruleBlder.arguments,
	}, nil
}

func NewRuleBuilder() *RuleBuilder {
	return &RuleBuilder{}
}
