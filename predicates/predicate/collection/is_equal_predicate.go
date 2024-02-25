// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package collection

import (
	"fmt"

	"github.com/conformize/conformize/common/typed"
	"github.com/conformize/conformize/internal/providers/api/schema"
	"github.com/conformize/conformize/internal/providers/api/schema/attributes"
	"github.com/conformize/conformize/predicates"
	"github.com/conformize/conformize/predicates/condition"
)

type IsEqualPredicate struct {
	PredicateBuilder predicates.PredicateBuilder
}

func (isEqPrd *IsEqualPredicate) Test(value typed.Valuable, args typed.Valuable) (bool, error) {
	if args == nil || value == nil {
		return false, fmt.Errorf("arguments or value cannot be nil")
	}

	elemVal, ok := value.(typed.Elementable)
	if !ok {
		return false, fmt.Errorf("expected a list value, got %s", value.Type().Name())
	}

	oElemVal, ok := args.(typed.Elementable)
	if !ok {
		return false, fmt.Errorf("expected a list value as argument, got %s", args.Type().Name())
	}

	if elemVal.Length() != oElemVal.Length() {
		return false, nil
	}

	if elemVal.Length() == 0 && oElemVal.Length() == 0 {
		return true, nil
	}

	switch val := elemVal.(type) {
	case *typed.ListValue:
		oListVal, ok := oElemVal.(*typed.ListValue)
		if !ok {
			return false, fmt.Errorf("expected a list value as argument")
		}
		if val.ElementsType.Hint().TypeHint() != oListVal.ElementsType.Hint().TypeHint() {
			return false, fmt.Errorf("cannot compare lists with elements of different type")
		}
	case *typed.TupleValue:
		oTupleVal, ok := oElemVal.(*typed.TupleValue)
		if !ok {
			return false, fmt.Errorf("expected a tuple value as argument")
		}

		if len(val.ElementsTypes) != len(oTupleVal.ElementsTypes) {
			return false, nil
		}

		for i, elemType := range val.ElementsTypes {
			if elemType.Hint() != oTupleVal.ElementsTypes[i].Hint() {
				return false, nil
			}
		}
	default:
		return false, fmt.Errorf("value must be a list or tuple")
	}

	valElements := elemVal.Items()
	oValElements := oElemVal.Items()

	startIdx, endIdx := 0, elemVal.Length()-1
	for startIdx <= endIdx {
		elemEqPrd, err := isEqPrd.PredicateBuilder.Build(condition.EQ)
		if err != nil {
			return false, err
		}

		ok, err := elemEqPrd.Test(valElements[startIdx], oValElements[startIdx])
		if !ok || err != nil {
			return false, err
		}

		if startIdx == endIdx {
			break
		}

		ok, err = elemEqPrd.Test(valElements[endIdx], oValElements[endIdx])
		if !ok || err != nil {
			return false, err
		}

		startIdx++
		endIdx--
	}
	return true, nil
}

func (listEqPrd *IsEqualPredicate) Schema() schema.Schemable {
	return &schema.Schema{
		Description: "List equality predicate",
		Version:     1,
		Attributes: map[string]schema.Attributeable{
			"Value": &attributes.ListAttribute{
				Required:     true,
				ElementsType: &typed.GenericTyped{},
			},
			"Arguments": &attributes.TupleAttribute{
				Required: true,
				ElementsTypes: []typed.Typeable{
					&typed.ListTyped{ElementsType: &typed.GenericTyped{}},
				},
			},
		},
	}
}
