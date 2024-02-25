// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package pathparser

import (
	"fmt"
	"strconv"

	"github.com/conformize/conformize/common/path"
)

type PathParser struct{}

func NewPathParser() *PathParser {
	return &PathParser{}
}

func (pParser *PathParser) Parse(pathStr string) (path.Steps, error) {
	if len(pathStr) == 0 {
		return path.Steps{}, nil
	}
	lexer := newLexer(pathStr)
	var pathSteps path.Steps
	for token := lexer.NextItem(); token != nil && token.itemType != EOL; token = lexer.NextItem() {
		switch token.itemType {
		case OBJECT_IDENTIFIER:
			pathSteps.Add(path.ObjectStep(token.value))
		case IDENTIFIER:
			pathSteps.Add(path.KeyStep(token.value))
		case PROPERTY:
			if token.value != LENGTH_PROPERTY {
				continue
			}
			pathSteps.Add(path.PropertyStep(token.value))
		case PROPERTY_NAME:
			pathSteps.Add(path.AttributeStep(token.value))
		case PROPERTY_INDEX:
			idx, err := strconv.ParseInt(token.value, 10, 64)
			if err != nil {
				return nil, err
			}
			pathSteps.Add(path.IndexStep(idx))
		case UNEXPECTED:
			return nil, fmt.Errorf("unexpected token %s, at position %d in %s", token.value, token.startPos+1, pathStr)
		}
	}
	return pathSteps, nil
}
