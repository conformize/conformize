// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package path

import (
	"fmt"

	"github.com/conformize/conformize/common/typed"
)

type ValuePathWalk struct {
	Destination       typed.Valuable
	CreateValueAtPath bool
}

func (vpw *ValuePathWalk) Walk(p *Path) (typed.Valuable, error) {
	if p == nil || len(p.Steps()) == 0 {
		return vpw.Destination, nil
	}

	steps := p.Steps()
	next, _ := steps.Next()

	var dstValues map[string]typed.Valuable
	var dstValType typed.Typeable

	stepName := next.String()
	switch dst := vpw.Destination.(type) {
	case *typed.ObjectValue:
		dstValues = dst.Fields
		dstValType = dst.FieldsTypes[stepName]
	case *typed.MapValue:
		dstValues = dst.Elements
		dstValType = dst.ElementsType
	default:
		return nil, fmt.Errorf("cannot walk '%s' type", dst.Type().Name())
	}

	if newDst, exists := dstValues[stepName]; exists {
		vpw.Destination = newDst
		return vpw.Walk(NewPath(steps))
	}

	if !vpw.CreateValueAtPath {
		return nil, fmt.Errorf("could not walk step - key or field '%s' not found", stepName)
	}

	if dstValType == nil {
		return nil, fmt.Errorf("cannot create value at path '%s' - type is not defined", stepName)
	}

	var err error
	if dstValues[stepName], err = typed.CreateValue(dstValType); err != nil {
		return nil, err
	}
	vpw.Destination = dstValues[stepName]
	return vpw.Walk(NewPath(steps))
}
