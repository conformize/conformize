// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package primitive

import (
	"regexp"

	"github.com/conformize/conformize/common/typed"
	"github.com/conformize/conformize/internal/providers/api/schema"
	"github.com/conformize/conformize/internal/providers/api/schema/attributes"
	"github.com/conformize/conformize/predicates/predicate"
)

type StringMatchesExpressionPredicate struct {
	predicate.PredicateArgumentsValidator
	Args *typed.TupleValue
}

func (strMatchExprPrd *StringMatchesExpressionPredicate) Test(value typed.Valuable) (bool, error) {
	valid, validErr := strMatchExprPrd.Validate(value, strMatchExprPrd.Args, strMatchExprPrd.Schema())
	if !valid {
		return valid, validErr
	}

	var s string
	value.As(&s)

	var expr string
	strMatchExprPrd.Args.Elements[0].As(&expr)

	var regExp, regExpErr = regexp.Compile(expr)
	if regExpErr != nil {
		return false, regExpErr
	}
	return regExp.MatchString(s), nil
}

func (strMatchExprPrd *StringMatchesExpressionPredicate) ArgumentsLength() int {
	return 1
}

func (strMatchExprPrd *StringMatchesExpressionPredicate) Arguments(args *typed.TupleValue) {
	strMatchExprPrd.Args = args
}

func (strMatchExprPrd *StringMatchesExpressionPredicate) Schema() schema.Schemable {
	return &schema.Schema{
		Description: "String regular expression match predicate",
		Version:     1,
		Attributes: map[string]schema.Attributeable{
			"Value": &attributes.StringAttribute{
				Required: true,
			},
			"Arguments": &attributes.TupleAttribute{
				Required: true,
				ElementsTypes: []typed.Typeable{
					&typed.StringTyped{},
				},
			},
		},
	}
}
