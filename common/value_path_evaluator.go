// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package common

import (
	"fmt"
	"reflect"

	"github.com/conformize/conformize/common/ds"
	"github.com/conformize/conformize/common/path"
	"github.com/conformize/conformize/common/reflected"
	"github.com/conformize/conformize/common/typed"
)

type ValuePathEvaluator struct{}

func (vpeval *ValuePathEvaluator) Evaluate(node *ds.Node[string, interface{}], p *path.Path) (*ds.Node[string, interface{}], error) {
	if node == nil || len(p.Steps()) == 0 {
		return node, nil
	}

	var current = node
	var children ds.NodeList[string, interface{}]
	var ok bool
	steps := p.Steps()
	nextStep, _ := steps.Next()

walkLoop:
	for nextStep != nil {
		switch nextStep := nextStep.(type) {
		case path.ObjectStep, path.KeyStep:
			children, ok = current.GetChildren(nextStep.String())
			if !ok {
				break walkLoop
			}
			current = children.First()
		case path.AttributeStep:
			if attr, ok := current.GetAttribute(nextStep.String()); ok {
				return attr, nil
			}
		case path.PropertyStep:
			switch nextStep.String() {
			case "length":
				reflectNodeRefVal := reflect.ValueOf(current.Value)
				vTypeHint := typed.TypeHintOf(reflectNodeRefVal)
				v, err := reflected.ValueFromTypeHint(reflectNodeRefVal, vTypeHint)
				if err != nil {
					return nil, err
				}

				val, ok := v.(typed.Lengthable)
				if !ok {
					return nil, fmt.Errorf("value of type %s has no property %s", v.Type().Name(), nextStep.String())
				}
				vNode := ds.NewNode[string, interface{}]()
				vNode.Value = val.Length()
				return vNode, nil
			default:
				return nil, fmt.Errorf("unknown property %s", nextStep.String())
			}
		case path.IndexStep:
			if children == nil {
				return nil, fmt.Errorf("index step '%s' cannot be applied without a preceding key", nextStep.String())
			}

			idx := int(nextStep)
			idx--
			if idx < 0 || idx >= children.Count() {
				return nil, fmt.Errorf("index %d out of range [1, %d)", idx+1, children.Count())
			}

			current = children.Get(idx)
		default:
			return nil, fmt.Errorf("invalid step in path")
		}
		nextStep, _ = steps.Next()
	}

	if current != nil {
		return current, nil
	}
	return nil, fmt.Errorf("step '%s' not found", nextStep)
}
