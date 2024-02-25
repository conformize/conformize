// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package common

import (
	"fmt"

	"github.com/conformize/conformize/common/ds"
	"github.com/conformize/conformize/common/path"
	"github.com/conformize/conformize/common/reflected"
	"github.com/conformize/conformize/common/typed"
)

type callbackFunc func(string, *ds.Node[string, interface{}])

func noOpCallbackFn(string, *ds.Node[string, interface{}]) {}

type ValuePathEvaluator struct{}

func (vpeval *ValuePathEvaluator) Evaluate(node *ds.Node[string, interface{}], p *path.Path, callback callbackFunc) (*ds.Node[string, interface{}], error) {
	if node == nil || len(p.Steps()) == 0 {
		return node, nil
	}

	nodeRef := node
	callbackFn := noOpCallbackFn
	if callback != nil {
		callbackFn = callback
	}

	steps := p.Steps()
	nextStep, _ := steps.Next()

walkLoop:
	for nextStep != nil {
		switch nextStep := nextStep.(type) {
		case path.ObjectStep, path.KeyStep:
			childRef, ok := nodeRef.GetChild(nextStep.String())
			nodeRef = childRef
			if !ok {
				break walkLoop
			}
			callbackFn(nextStep.String(), childRef)
		case path.AttributeStep:
			if attr, ok := nodeRef.GetAttribute(nextStep.String()); ok {
				callbackFn(nextStep.String(), attr)
				return attr, nil
			}
		case path.PropertyStep:
			switch nextStep.String() {
			case "length":
				vTypeHint := typed.TypeHintOf(nodeRef.Value)
				v, err := reflected.ValueFromTypeHint(nodeRef.Value, vTypeHint)
				if err != nil {
					return nil, err
				}

				val, ok := v.(typed.Lengthable)
				if !ok {
					return nil, fmt.Errorf("value of type %s has no property %s", v.Type().Name(), nextStep.String())
				}
				vNode := ds.NewNode[string, interface{}]()
				vNode.Value = val.Length()
				callbackFn(nextStep.String(), vNode)
				return vNode, nil
			default:
				return nil, fmt.Errorf("unknown property %s", nextStep.String())
			}
		case path.IndexStep:
			idx := int(nextStep)
			for nodeIdx := 0; nodeIdx <= idx-1; nodeIdx++ {
				nodeRef = nodeRef.Next()
				if nodeRef == nil {
					return nil, fmt.Errorf("index %d out out of range [1, %d]", idx, nodeIdx+1)
				}
				callbackFn(nextStep.String(), nodeRef)
			}
		default:
			return nil, fmt.Errorf("invalid step in path")
		}
		nextStep, _ = steps.Next()
	}

	if nodeRef != nil {
		return nodeRef, nil
	}
	return nil, fmt.Errorf("step '%s' not found", nextStep)
}
