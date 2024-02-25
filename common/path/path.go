// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package path

import (
	"fmt"
	"strings"
)

const defaultSeparator rune = '.'

type Path struct {
	steps Steps
}

func (p *Path) Steps() Steps {
	if len(p.steps) == 0 {
		return Steps{}
	}
	return p.steps.Clone()
}

func (p *Path) Clone() Path {
	return Path{
		steps: p.steps.Clone(),
	}
}

func NewPath(steps Steps) *Path {
	return &Path{
		steps: steps,
	}
}

func NewFromString(path string) (*Path, error) {
	return newFromString(path, defaultSeparator)
}

func NewFromStringWithSeparator(path string, sep rune) (*Path, error) {
	return newFromString(path, sep)
}

func (p *Path) String() string {
	if len(p.steps) == 0 {
		return ""
	}
	return p.steps.String()
}

func newFromString(path string, sep rune) (*Path, error) {
	steps := strings.FieldsFunc(path, func(r rune) bool {
		return r == sep
	})

	stepsCount := len(steps)
	if stepsCount == 0 {
		return nil, fmt.Errorf("malformed path: %s", path)
	}
	pathSteps := Steps{}
	for _, step := range steps {
		pathSteps.Add(KeyStep(step))
	}
	return &Path{
		steps: pathSteps,
	}, nil
}
