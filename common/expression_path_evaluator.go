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
	"strconv"

	"github.com/conformize/conformize/common/ds"
	"github.com/conformize/conformize/common/path"
	"github.com/conformize/conformize/common/reflected"
	"github.com/conformize/conformize/common/typed"
	"github.com/conformize/conformize/common/typed/functions"
	"github.com/conformize/conformize/predicates"
)

type ExpressionPathEvaluator struct{}

type IterFnNodeValue struct {
	Fn   func(it typed.Iterable, p predicates.Predicate, args typed.Valuable) (bool, error)
	Iter typed.Iterable
}

func (vpeval *ExpressionPathEvaluator) Evaluate(node *ds.Node[string, any], p *path.Path) (*ds.Node[string, any], error) {
	if node == nil || len(p.Steps()) == 0 {
		return node, nil
	}

	var current = node
	var children ds.NodeList[string, any]
	var ok bool
	steps := p.Steps()
	nextStep, _ := steps.Next()

	var err error
	var v typed.Valuable
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
				v, err = reflected.ValueFromTypeHint(reflectNodeRefVal, vTypeHint)
				if err != nil {
					return nil, err
				}

				val, ok := v.(typed.Lengthable)
				if !ok {
					return nil, fmt.Errorf("value of type %s has no property %s", v.Type().Name(), nextStep.String())
				}
				vNode := ds.NewNode[string, any]()
				vNode.Key = nextStep.String()
				vNode.Value = val.Length()
				return vNode, nil
			default:
				return nil, fmt.Errorf("unknown property %s", nextStep.String())
			}
		case path.FunctionStep:
			vNode := ds.NewNode[string, any]()
			vNode.Key = nextStep.String()

			reflectNodeRefVal := reflect.ValueOf(current.Value)
			vTypeHint := typed.TypeHintOf(reflectNodeRefVal)
			v, err = reflected.ValueFromTypeHint(reflectNodeRefVal, vTypeHint)
			if err != nil {
				return nil, err
			}

			val, ok := v.(typed.Elementable)
			if !ok {
				return nil, fmt.Errorf("value of type %s is not iterable", v.Type().Name())
			}

			switch nextStep.String() {
			case "none":
				vNode.Value = &IterFnNodeValue{
					Fn:   functions.NoneOf,
					Iter: typed.NewElementIterator(val),
				}
			case "any":
				vNode.Value = &IterFnNodeValue{
					Fn:   functions.AnyOf,
					Iter: typed.NewElementIterator(val),
				}
			case "each":
				vNode.Value = &IterFnNodeValue{
					Fn:   functions.Each,
					Iter: typed.NewElementIterator(val),
				}
			default:
				return nil, fmt.Errorf("unknown function '%s'", nextStep.String())
			}
			return vNode, nil
		case path.IndexStep:
			if children == nil {
				return nil, fmt.Errorf("index step '%s' cannot be applied without a preceding key", nextStep.String())
			}

			idx, err := strconv.Atoi(nextStep.String())
			if err != nil {
				return nil, fmt.Errorf("invalid index '%s': %w", nextStep.String(), err)
			}

			if current.Value != nil {
				reflectVal := reflect.ValueOf(current.Value)
				if reflectVal.Type().Kind() == reflect.Slice &&
					!reflectVal.IsNil() && typed.TypeHintOf(reflectVal).TypeHint() == typed.List {

					if idx < 0 || idx >= reflectVal.Len() {
						return nil, fmt.Errorf("index %d out of range [0, %d)", idx+1, reflectVal.Len())
					}
					vNode := ds.NewNode[string, any]()
					vNode.Key = nextStep.String()
					vNode.Value = reflectVal.Index(idx).Interface()
					return vNode, nil
				}
			}

			if idx < 0 || idx >= children.Count() {
				return nil, fmt.Errorf("index %d out of range [0, %d)", idx+1, children.Count())
			}

			current = children.Get(idx)
		case *path.AttributeStep:
			if attr, ok := current.GetAttribute(nextStep.String()); ok {
				return attr, nil
			}
			return nil, fmt.Errorf("attribute '%s' not found in node '%s'", nextStep.String(), current.Key)
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
