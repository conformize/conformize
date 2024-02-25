// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package pathparser

import (
	"fmt"
)

type tokenItem struct {
	itemType tokenType
	value    string
	startPos int
}

type lexer struct {
	input               string
	pos                 int
	scanBuf             []byte
	currentState        *lexState
	lastPos             int
	tokenItems          []tokenItem
	lexStateTransitions *lexStateTransitionMap
}

type lexState struct {
	state     tokenType
	accepts   lexStateAcceptFn
	nextState lexStateTransitionFn
}

func (l *lexer) skip() {
	l.pos++
}

func (l *lexer) accepts() bool {
	return l.currentState.accepts(l)
}

func (l *lexer) NextItem() *tokenItem {
	tokensCount := len(l.tokenItems)
	if tokensCount == 0 {
		return nil
	}
	item := l.tokenItems[0]
	if tokensCount > 1 {
		l.tokenItems = l.tokenItems[1:]
	} else {
		l.tokenItems = []tokenItem{}
	}
	return &item
}

func newLexer(input string) *lexer {
	l := &lexer{input: input, lexStateTransitions: stateTransitions()}
	l.currentState = (*l.lexStateTransitions)[OBJECT_IDENTIFIER_START]
	l.run()
	return l
}

func (l *lexer) add(tType tokenType) {
	token := tokenItem{itemType: tType, value: string(l.scanBuf), startPos: l.lastPos}
	l.tokenItems = append(l.tokenItems, token)
	l.lastPos += len(l.scanBuf)
	l.scanBuf = []byte{}
}

func (l *lexer) run() {
	if l.currentState.state != OBJECT_IDENTIFIER_START || !l.accepts() {
		l.scanBuf = append(l.scanBuf, l.input[l.pos])
		l.add(UNEXPECTED)
		return
	}
	l.scanBuf = append(l.scanBuf, l.input[l.pos])
	l.add(l.currentState.state)
	l.skip()
	l.currentState = (*l.lexStateTransitions)[l.currentState.nextState(l)]
	for {
		if l.accepts() {
			l.scanBuf = append(l.scanBuf, l.input[l.pos])
			l.skip()
			continue
		}

		if l.currentState.state == KEYWORD {
			curr := string(l.scanBuf)
			if isKeyword(curr) {
				switch curr {
				case ATTRIBUTES_PROP:
					l.add(PROPERTY)
				case LENGTH_PROP:
					l.add(PROPERTY)
				case NONE_FN, ANY_FN, EACH_FN:
					l.add(FUNCTION)
				}
			}
		} else {
			l.add(l.currentState.state)
		}

		nextState := (*l.lexStateTransitions)[l.currentState.nextState(l)]
		if nextState.state == EOL || l.currentState.state == UNEXPECTED {
			l.add(nextState.state)
			break
		}
		l.currentState = nextState

	}
}

func (l *lexer) peek() (*byte, error) {
	if l.pos >= len(l.input) {
		return nil, fmt.Errorf("end of input")
	}
	ch := l.input[l.pos]
	return &ch, nil
}
