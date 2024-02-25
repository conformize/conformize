// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package primitive

import (
	"fmt"
	"regexp"

	"github.com/conformize/conformize/common/typed"
	"github.com/conformize/conformize/internal/providers/api/schema"
	"github.com/conformize/conformize/internal/providers/api/schema/attributes"
)

type StringMatchesExpressionPredicate struct {
	Args typed.Valuable
}

func (strMatchExprPrd *StringMatchesExpressionPredicate) Test(value typed.Valuable) (bool, error) {
	if value == nil {
		return false, fmt.Errorf("value is nil")
	}

	if value.Type().Hint() != typed.String {
		return false, fmt.Errorf("expected a string value, got %s", value.Type().Name())
	}

	var s string
	value.As(&s)

	if strMatchExprPrd.Args == nil {
		return false, fmt.Errorf("arguments is nil")
	}

	if strMatchExprPrd.Args.Type().Hint() != typed.String {
		return false, fmt.Errorf("expected a string argument, got %s", strMatchExprPrd.Args.Type().Name())
	}

	var expr string
	strMatchExprPrd.Args.As(&expr)

	var regExp, regExpErr = regexp.Compile(expr)
	if regExpErr != nil {
		return false, regExpErr
	}
	return regExp.MatchString(s), nil
}

func (strMatchExprPrd *StringMatchesExpressionPredicate) ArgumentsCount() int {
	return 1
}

func (strMatchExprPrd *StringMatchesExpressionPredicate) Arguments(args typed.Valuable) {
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
