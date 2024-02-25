// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package pathparser

import (
	"fmt"

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
			switch token.value {
			case ATTRIBUTES_PROP:
				var next *tokenItem
				for _, expected := range []tokenType{DELIMITER, QUOTE} {
					if next = lexer.NextItem(); next == nil {
						return nil, fmt.Errorf("expected token after '%s' at position %d in %s", token.value, token.startPos, pathStr)
					}
					if next.itemType != expected {
						return nil, fmt.Errorf("unexpected token after '%s' at position %d in %s", token.value, token.startPos, pathStr)
					}
				}

				next = lexer.NextItem()
				if next == nil {
					return nil, fmt.Errorf("expected attribute name token after '%s' at position %d in %s", token.value, token.startPos, pathStr)
				}

				if next.itemType != IDENTIFIER {
					return nil, fmt.Errorf("unexpected token after '%s' at position %d in %s", token.value, token.startPos, pathStr)
				}
				pathSteps.Add(path.AttributeStep(next.value))
			case LENGTH_PROP:
				pathSteps.Add(path.PropertyStep(token.value))
			default:
				return nil, fmt.Errorf("unexpected property %s, at position %d in %s", token.value, token.startPos, pathStr)
			}
		case FUNCTION:
			switch token.value {
			case EACH_FN, NONE_FN, ANY_FN:
				pathSteps.Add(path.FunctionStep(token.value))
			default:
				return nil, fmt.Errorf("unexpected function %s, at position %d in %s", token.value, token.startPos, pathStr)
			}
		case INDEX:
			pathSteps.Add(path.IndexStep(token.value))
		case DELIMITER, QUOTE, OBJECT_IDENTIFIER_START:
			continue
		case UNEXPECTED:
			return nil, fmt.Errorf("unexpected token %s, at position %d in %s", token.value, token.startPos, pathStr)
		case EOL:
			return pathSteps, nil
		}
	}
	return pathSteps, nil
}
