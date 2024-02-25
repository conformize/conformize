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
	KEYWORD
	FUNCTION
	PROPERTY
	INDEX
	QUOTE
	DELIMITER
	EOL
)

const (
	ATTRIBUTES_PROP = "attributes"
	LENGTH_PROP     = "length"
	NONE_FN         = "none"
	ANY_FN          = "any"
	EACH_FN         = "each"
)

type lexStateAcceptFn func(l *lexer) bool
type lexStateTransitionFn func(l *lexer) tokenType

type lexStateTransitionMap map[tokenType]*lexState

func isIdentifierCharacter(ch *byte) bool {
	return *ch != '\''
}

func isKeywordCharacter(ch *byte) bool {
	return *ch >= 'a' && *ch <= 'z' || *ch >= 'A' && *ch <= 'Z'
}

func isObjectIdentifierCharacter(ch *byte) bool {
	return *ch >= 'a' && *ch <= 'z' || *ch >= 'A' && *ch <= 'Z' || *ch == '_' || *ch == '-'
}

func isIndexCharacter(ch *byte) bool {
	return *ch >= '0' && *ch <= '9'
}

var keywordsExp = `\b(attributes|length|each|any|none)\b`
var keywordsRegexp = regexp.MustCompile(keywordsExp)

func isKeyword(s string) bool {
	return keywordsRegexp.MatchString(s)
}

func lexStateTransitions() *lexStateTransitionMap {
	return &lexStateTransitionMap{
		OBJECT_IDENTIFIER_START: &lexState{
			state: OBJECT_IDENTIFIER_START,
			accepts: func(l *lexer) bool {
				return l.input[l.pos] == '$'
			},
			nextState: func(l *lexer) tokenType {
				ch, err := l.peek()
				if err != nil {
					return EOL
				}

				if isObjectIdentifierCharacter(ch) {
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
				ch, err := l.peek()
				if err != nil {
					return EOL
				}

				if *ch == '.' {
					return DELIMITER
				}
				return UNEXPECTED
			},
		},
		DELIMITER: &lexState{
			state: DELIMITER,
			accepts: func(l *lexer) bool {
				ch, err := l.peek()
				return ch != nil && err == nil && *ch == '.'
			},
			nextState: func(l *lexer) tokenType {
				ch, err := l.peek()
				if err != nil {
					return EOL
				}

				if *ch == '\'' {
					return QUOTE
				}

				if isIndexCharacter(ch) {
					return INDEX
				}

				if isKeywordCharacter(ch) {
					return KEYWORD
				}

				return UNEXPECTED
			},
		},
		QUOTE: &lexState{
			state: QUOTE,
			accepts: func(l *lexer) bool {
				ch, err := l.peek()
				return ch != nil && err == nil && *ch == '\''
			},
			nextState: func(l *lexer) tokenType {
				ch, err := l.peek()
				if err != nil {
					return EOL
				}

				if *ch == '.' {
					return DELIMITER
				}

				if isIdentifierCharacter(ch) {
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
				ch, err := l.peek()
				if err != nil {
					return EOL
				}

				if *ch == '\'' {
					return QUOTE
				}

				if *ch == '.' {
					return DELIMITER
				}

				return UNEXPECTED
			},
		},
		INDEX: &lexState{
			state: INDEX,
			accepts: func(l *lexer) bool {
				ch, err := l.peek()
				return ch != nil && err == nil && isIndexCharacter(ch)
			},
			nextState: func(l *lexer) tokenType {
				ch, err := l.peek()
				if err != nil {
					return EOL
				}

				if *ch == '\'' {
					return QUOTE
				}

				if *ch == '.' {
					return DELIMITER
				}
				return UNEXPECTED
			},
		},
		KEYWORD: &lexState{
			state: KEYWORD,
			accepts: func(l *lexer) bool {
				ch, err := l.peek()
				return ch != nil && err == nil && isKeywordCharacter(ch)
			},
			nextState: func(l *lexer) tokenType {
				ch, err := l.peek()
				if err != nil {
					return EOL
				}

				if *ch == '\'' {
					return QUOTE
				}

				if *ch == '.' {
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
