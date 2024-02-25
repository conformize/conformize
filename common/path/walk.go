// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package path

import (
	"fmt"

	"github.com/conformize/conformize/common/reflected"
	"github.com/conformize/conformize/common/typed"
)

type Walk struct {
	Destination       typed.Valuable
	CreateValueAtPath bool
}

func (w *Walk) Walk(p *Path) (typed.Valuable, error) {
	if p == nil || len(p.Steps()) == 0 {
		return w.Destination, nil
	}

	steps := p.Steps()
	next, _ := steps.Next()

	var createValType typed.Typeable
	var createNewDst map[string]typed.Valuable

	stepName := next.String()
	switch dst := w.Destination.(type) {
	case *typed.ObjectValue:
		dstField, exists := dst.Fields[stepName]
		if exists {
			w.Destination = dstField
		} else if w.CreateValueAtPath {
			createNewDst = dst.Fields
			createValType = dst.FieldsTypes[stepName]
		} else {
			return nil, fmt.Errorf("could not walk step - field '%s' not found", stepName)
		}
	case *typed.MapValue:
		dstElem, exists := dst.Elements[stepName]
		if exists {
			w.Destination = dstElem
		} else if w.CreateValueAtPath {
			createNewDst = dst.Elements
			createValType = dst.ElementsType
		} else {
			return nil, fmt.Errorf("could not walk step - key '%s' not found", stepName)
		}
	default:
		return nil, fmt.Errorf("cannot walk '%s' type", dst.Type().Name())
	}

	if createNewDst != nil {
		var err error
		createNewDst[stepName], err = reflected.Value(nil, createValType)
		if err != nil {
			return nil, err
		}
		w.Destination = createNewDst[stepName]
	}
	return w.Walk(NewPath(steps))
}
