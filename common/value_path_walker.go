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
)

type ValuePathWalker struct{}

func (vpw *ValuePathWalker) Walk(node *ds.Node[string, interface{}], p *path.Path) (*ds.Node[string, interface{}], error) {
	if node == nil || len(p.Steps()) == 0 {
		return node, nil
	}

	steps := p.Steps()
	nextStep, _ := steps.Next()
	switch nextStep := nextStep.(type) {
	case path.ObjectStep, path.KeyStep:
		if child, ok := node.GetChild(nextStep.String()); ok {
			return vpw.Walk(child, path.NewPath(steps))
		}
	case path.AttributeStep:
		if attr, ok := node.GetAttribute(nextStep.String()); ok {
			return attr, nil
		}
	case path.IndexStep:
		idx := int(nextStep)
		for nodeIdx := 0; nodeIdx < idx; nodeIdx++ {
			node = node.Next()
			if node == nil {
				return nil, fmt.Errorf("index %d out of range", idx)
			}
		}
		return vpw.Walk(node, path.NewPath(steps))
	default:
		return nil, fmt.Errorf("invalid step in path")
	}
	return nil, fmt.Errorf("step %s not found", nextStep)
}
