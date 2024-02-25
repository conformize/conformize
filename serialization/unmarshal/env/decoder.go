// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package env

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode/utf16"
	"unicode/utf8"

	"github.com/conformize/conformize/serialization/unmarshal/functions"
)

type Decoder struct {
	bufReader *bufio.Reader
}

var errInvalidKey = errors.New("invalid variable name")
var errEOF = errors.New("end of file")
var errInvalidLineFormat = errors.New("invalid line format. No key-value pair found")

func findSeparatorPosition(line string) int {
	return strings.Index(line, "=")
}

func isEmpty(line string) bool {
	return len(line) == 0
}

func isNewline(line string) bool {
	return line[0] == '\n' || line[0] == '\r'
}

func isComment(line string) bool {
	return line[0] == '#'
}

func isMultilineToken(line string) bool {
	return len(line) == 3 && (line == `"""` || line == `'''`)
}

func isValidKey(key *string) bool {
	if len(*key) == 0 {
		return false
	}

	// First character must be letter or underscore (POSIX standard)
	firstChar := (*key)[0]
	if !(isLetter(firstChar) || isUnderscore(firstChar)) {
		return false
	}

	// Subsequent characters can be letters, digits, or underscores
	for pos := 1; pos < len(*key); pos++ {
		char := (*key)[pos]
		if !(isLetter(char) || isDigit(char) || isUnderscore(char)) {
			return false
		}
	}
	return true
}

func isLetter(char byte) bool {
	return isLowercase(char) || isUppercase(char)
}

func isLowercase(char byte) bool {
	return 97 <= char && char <= 122
}

func isDigit(char byte) bool {
	return 48 <= char && char <= 57
}

func isUppercase(char byte) bool {
	return 65 <= char && char <= 90
}

func isUnderscore(char byte) bool {
	return char == 95
}

func extractKey(line string, endPos int) (*string, error) {
	key := strings.TrimSpace(line[:endPos])
	if isValidKey(&key) {
		return &key, nil
	}
	return nil, errInvalidKey
}

func extractValue(line string, startPos int) string {
	return strings.TrimSpace(line[startPos:])
}

func parseHexString(hexStr string) (uint16, error) {
	if len(hexStr) != 4 {
		return 0, fmt.Errorf("invalid unicode escape: expected 4 hex digits, got %d", len(hexStr))
	}

	val, err := strconv.ParseUint(hexStr, 16, 16)
	if err != nil {
		return 0, fmt.Errorf("invalid unicode escape sequence '\\u%s': %w", hexStr, err)
	}
	return uint16(val), nil
}

func parseValue(value string) (any, error) {
	if val, err := functions.DecodeStringValue(value); err == nil {
		return val, nil
	}
	return parseString(value)
}

func decodeUnicode(val string) (rune, error) {
	if len(val) < 4 {
		return 0, fmt.Errorf("incomplete unicode escape sequence: expected 4 hex digits")
	}

	hexValue, err := parseHexString(val[0:4])
	if err != nil {
		return 0, err
	}

	decoded := utf16.Decode([]uint16{hexValue})
	if len(decoded) == 0 {
		return 0, fmt.Errorf("invalid unicode code point")
	}
	return decoded[0], nil
}

func escapeCharacter(char byte) rune {
	switch char {
	case 't':
		return '\t'
	case 'r':
		return '\r'
	case 'n':
		return '\n'
	case 'f':
		return '\f'
	case 'b':
		return '\b'
	default:
		return rune(char)
	}
}

func parseString(value string) (string, error) {
	if len(value) == 0 {
		return "", nil
	}

	// Handle quoted strings
	if (value[0] == '"' || value[0] == '\'') && len(value) >= 2 {
		quote := value[0]
		if value[len(value)-1] == quote {
			// Remove surrounding quotes and parse the content
			return parseQuotedString(value[1:len(value)-1], quote == '"')
		}
		return "", fmt.Errorf("unterminated quoted string")
	}

	// Handle unquoted strings (stop at # for inline comments)
	return parseUnquotedString(value)
}

func parseQuotedString(value string, allowEscapes bool) (string, error) {
	if !allowEscapes {
		// Single quotes - literal string, no escape processing
		return value, nil
	}

	// Double quotes - process escape sequences
	buf := bytes.Buffer{}
	valLen := len(value)
	offset := 0

	for offset < valLen {
		char := value[offset]
		if char == '\\' {
			nextPos := offset + 1
			if nextPos < valLen {
				offset = nextPos
				char = value[offset]
				if char == 'u' {
					offset++
					if offset+4 > valLen {
						return "", fmt.Errorf("incomplete unicode escape at position %d", offset-2)
					}
					hexPart := value[offset : offset+4]
					decoded, err := decodeUnicode(hexPart)
					if err != nil {
						return "", fmt.Errorf("invalid unicode escape at position %d: %w", offset-2, err)
					}
					buf.WriteRune(decoded)
					offset += 4
					continue
				} else {
					buf.WriteRune(escapeCharacter(char))
				}
			} else {
				// Trailing backslash - treat as literal
				buf.WriteRune('\\')
			}
			offset++
		} else {
			decoded, decodedLen := utf8.DecodeRuneInString(value[offset:])
			if decoded == utf8.RuneError && decodedLen == 1 {
				return "", fmt.Errorf("invalid UTF-8 sequence at position %d", offset)
			}
			buf.WriteRune(decoded)
			offset += decodedLen
		}
	}
	return buf.String(), nil
}

func parseUnquotedString(value string) (string, error) {
	// For unquoted strings, stop at first # (inline comment)
	if commentPos := strings.Index(value, "#"); commentPos != -1 {
		value = value[:commentPos]
	}
	return strings.TrimSpace(value), nil
}

func (d *Decoder) Decode() (key *string, value any, err error) {
	line, err := d.bufReader.ReadString('\n')
	if err != nil {
		if err == io.EOF {
			return nil, nil, io.EOF
		}
		return nil, nil, fmt.Errorf("error reading line: %w", err)
	}

	lineTrim := strings.TrimSpace(line)
	if isEmpty(lineTrim) || isComment(lineTrim) || isNewline(lineTrim) {
		return nil, nil, nil
	}

	separatorPos := findSeparatorPosition(lineTrim)
	if separatorPos == -1 {
		return nil, nil, errInvalidLineFormat
	}

	key, err = extractKey(lineTrim, separatorPos)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid key in line '%s': %w", lineTrim, err)
	}

	separatorPos++
	valuePart := extractValue(lineTrim, separatorPos)
	if isMultilineToken(valuePart) {
		token := valuePart
		closed := false

		bldr := strings.Builder{}
		for {
			line, err = d.bufReader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					break
				}
				return nil, nil, fmt.Errorf("error reading multiline continuation: %w", err)
			}

			valuePart = strings.TrimSpace(line)
			if isMultilineToken(valuePart) {
				if token == valuePart {
					closed = true
					break
				}
				continue
			}

			if bldr.Len() > 0 {
				bldr.WriteByte('\n')
			}

			parsedLine, err := parseString(valuePart)
			if err != nil {
				return nil, nil, fmt.Errorf("error parsing multiline content '%s': %w", valuePart, err)
			}
			bldr.WriteString(parsedLine)
		}

		if !closed {
			return nil, nil, fmt.Errorf("malformed multiline value! Expected closing token %s not found", token)
		}
		return key, bldr.String(), nil
	}

	parsedValue, err := parseValue(valuePart)
	if err != nil {
		return nil, nil, fmt.Errorf("error parsing value '%s': %w", valuePart, err)
	}
	return key, parsedValue, nil
}

func NewDecoder(bufReader *bufio.Reader) *Decoder {
	return &Decoder{bufReader: bufReader}
}
