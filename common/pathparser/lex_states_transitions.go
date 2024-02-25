// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package pathparser

import (
	"regexp"
	"sync"
)

type tokenType int

const (
	UNEXPECTED tokenType = iota - 1
	OBJECT_IDENTIFIER_START
	OBJECT_IDENTIFIER
	IDENTIFIER
	PROPERTY
	PROPERTY_NAME
	PROPERTY_DELIMITER
	PROPERTY_OPENING_QUOTE
	PROPERTY_CLOSING_QUOTE
	PROPERTY_INDEX
	OPENING_QUOTE
	CLOSING_QUOTE
	DELIMITER
	EOL
)

const (
	LENGTH_PROPERTY = "length"
)

type lexStateAcceptFn func(l *lexer) bool
type lexStateTransitionFn func(l *lexer) tokenType

type lexStateTransitionMap map[tokenType]*lexState

func isIdentifierCharacter(ch *byte) bool {
	return *ch != '\''
}

func isPropertyNameCharacter(ch *byte) bool {
	return *ch >= 'a' && *ch <= 'z'
}

func isObjectIdentifierCharacter(ch *byte) bool {
	return *ch >= 'a' && *ch <= 'z' || *ch >= 'A' && *ch <= 'Z' || *ch == '_' || *ch == '-'
}

func isIndexCharacter(ch *byte) bool {
	return *ch >= '0' && *ch <= '9'
}

var keywordsExp = `\b(attributes|length)\b`
var keywordsRegexp = regexp.MustCompile(keywordsExp)

func isKeyword(s string) bool {
	return keywordsRegexp.MatchString(s)
}

func lexStateTransitions() *lexStateTransitionMap {
	return &lexStateTransitionMap{
		OBJECT_IDENTIFIER_START: &lexState{
			state: OBJECT_IDENTIFIER_START,
			accepts: func(l *lexer) bool {
				ch, err := l.peek()
				return l.pos == 0 && ch != nil && err == nil && *ch == '$'
			},
			nextState: func(l *lexer) tokenType {
				if ch, err := l.peek(); err != nil {
					return EOL
				} else if isObjectIdentifierCharacter(ch) {
					return OBJECT_IDENTIFIER
				}
				return UNEXPECTED
			},
		},
		OBJECT_IDENTIFIER: &lexState{
			state: OBJECT_IDENTIFIER,
			accepts: func(l *lexer) bool {
				ch, err := l.peek()
				return ch != nil && err == nil && isObjectIdentifierCharacter(ch)
			},
			nextState: func(l *lexer) tokenType {
				if ch, err := l.peek(); err != nil {
					return EOL
				} else if *ch == '.' {
					return DELIMITER
				}
				return UNEXPECTED
			},
		},
		DELIMITER: &lexState{
			state: DELIMITER,
			accepts: func(l *lexer) bool {
				ch, err := l.peek()
				return ch != nil && err == nil && l.lastPos == l.pos && *ch == '.'
			},
			nextState: func(l *lexer) tokenType {
				if ch, err := l.peek(); err != nil {
					return EOL
				} else if *ch == '\'' {
					return OPENING_QUOTE
				} else if isIndexCharacter(ch) {
					return PROPERTY_INDEX
				} else if isObjectIdentifierCharacter(ch) {
					return PROPERTY
				} else {
					return UNEXPECTED
				}
			},
		},
		OPENING_QUOTE: &lexState{
			state: OPENING_QUOTE,
			accepts: func(l *lexer) bool {
				ch, err := l.peek()
				return ch != nil && err == nil && l.lastPos == l.pos && *ch == '\''
			},
			nextState: func(l *lexer) tokenType {
				if ch, err := l.peek(); err != nil {
					return EOL
				} else if isIdentifierCharacter(ch) {
					return IDENTIFIER
				}
				return UNEXPECTED
			},
		},
		IDENTIFIER: &lexState{
			state: IDENTIFIER,
			accepts: func(l *lexer) bool {
				ch, err := l.peek()
				return ch != nil && err == nil && isIdentifierCharacter(ch)
			},
			nextState: func(l *lexer) tokenType {
				if ch, err := l.peek(); err != nil {
					return EOL
				} else if *ch == '\'' {
					return CLOSING_QUOTE
				}
				return UNEXPECTED
			},
		},
		CLOSING_QUOTE: &lexState{
			state: CLOSING_QUOTE,
			accepts: func(l *lexer) bool {
				ch, err := l.peek()
				return ch != nil && err == nil && l.lastPos == l.pos && *ch == '\''
			},
			nextState: func(l *lexer) tokenType {
				if ch, err := l.peek(); err != nil {
					return EOL
				} else if *ch == '.' {
					return DELIMITER
				}
				return UNEXPECTED
			},
		},
		PROPERTY: &lexState{
			state: PROPERTY,
			accepts: func(l *lexer) bool {
				ch, err := l.peek()
				return ch != nil && err == nil && isObjectIdentifierCharacter(ch)
			},
			nextState: func(l *lexer) tokenType {
				if ch, err := l.peek(); err != nil {
					return EOL
				} else if *ch == '.' {
					return PROPERTY_DELIMITER
				}
				return UNEXPECTED
			},
		},
		PROPERTY_DELIMITER: &lexState{
			state: PROPERTY_DELIMITER,
			accepts: func(l *lexer) bool {
				ch, err := l.peek()
				return ch != nil && err == nil && l.lastPos == l.pos && *ch == '.'
			},
			nextState: func(l *lexer) tokenType {
				if ch, err := l.peek(); err != nil {
					return EOL
				} else if *ch == '\'' {
					return PROPERTY_OPENING_QUOTE
				}
				return UNEXPECTED
			},
		},
		PROPERTY_OPENING_QUOTE: &lexState{
			state: PROPERTY_OPENING_QUOTE,
			accepts: func(l *lexer) bool {
				ch, err := l.peek()
				return ch != nil && err == nil && l.lastPos == l.pos && *ch == '\''
			},
			nextState: func(l *lexer) tokenType {
				if ch, err := l.peek(); err != nil {
					return EOL
				} else if isIdentifierCharacter(ch) {
					return PROPERTY_NAME
				}
				return UNEXPECTED
			},
		},
		PROPERTY_NAME: &lexState{
			state: PROPERTY_NAME,
			accepts: func(l *lexer) bool {
				ch, err := l.peek()
				return ch != nil && err == nil && isPropertyNameCharacter(ch)
			},
			nextState: func(l *lexer) tokenType {
				if ch, err := l.peek(); err != nil {
					return EOL
				} else if *ch == '\'' {
					return PROPERTY_CLOSING_QUOTE
				}
				return UNEXPECTED
			},
		},
		PROPERTY_CLOSING_QUOTE: &lexState{
			state: PROPERTY_CLOSING_QUOTE,
			accepts: func(l *lexer) bool {
				ch, err := l.peek()
				return ch != nil && err == nil && l.lastPos == l.pos && *ch == '\''
			},
			nextState: func(l *lexer) tokenType {
				if _, err := l.peek(); err != nil {
					return EOL
				}
				return UNEXPECTED
			},
		},
		PROPERTY_INDEX: &lexState{
			state: PROPERTY_INDEX,
			accepts: func(l *lexer) bool {
				ch, err := l.peek()
				return ch != nil && err == nil && isIndexCharacter(ch)
			},
			nextState: func(l *lexer) tokenType {
				if ch, err := l.peek(); err != nil {
					return EOL
				} else if *ch == '.' {
					return DELIMITER
				}
				return UNEXPECTED
			},
		},
		EOL: &lexState{
			state: EOL,
			accepts: func(l *lexer) bool {
				return false
			},
			nextState: func(l *lexer) tokenType {
				return UNEXPECTED
			},
		},
		UNEXPECTED: &lexState{
			state: UNEXPECTED,
			accepts: func(l *lexer) bool {
				return false
			},
			nextState: func(l *lexer) tokenType {
				return UNEXPECTED
			},
		},
	}
}

var (
	stateTransitionsMap *lexStateTransitionMap
	once                sync.Once
)

func stateTransitions() *lexStateTransitionMap {
	once.Do(func() {
		stateTransitionsMap = lexStateTransitions()
	})
	return stateTransitionsMap
}
