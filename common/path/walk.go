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
	var err error

	stepName := next.String()
	switch dst := w.Destination.(type) {
	case *typed.ObjectValue:
		if dstField, exists := dst.Fields[stepName]; exists {
			w.Destination = dstField
		} else if !w.CreateValueAtPath {
			err = fmt.Errorf("could not find field %s", stepName)
		} else {
			if dst.Fields[stepName], err = reflected.Value(nil, dst.FieldsTypes[stepName]); err == nil {
				w.Destination = dst.Fields[stepName]
			}
		}
	case *typed.MapValue:
		if dstElem, exists := dst.Elements[stepName]; exists {
			w.Destination = dstElem
		} else if !w.CreateValueAtPath {
			err = fmt.Errorf("could not find key %s", stepName)
		} else {
			if dst.Elements[stepName], err = reflected.Value(nil, dst.ElementsType); err == nil {
				w.Destination = dst.Elements[stepName]
			}
		}
	default:
		return nil, fmt.Errorf("type %s is not traversable", dst.Type().Name())
	}

	if err != nil {
		return nil, err
	}
	newPath := NewPath(steps)
	return w.Walk(newPath)
}
