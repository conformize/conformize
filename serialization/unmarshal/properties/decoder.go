// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package properties

import (
	"bufio"
	"bytes"
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

func findSeparatorPosition(line string) int {
	separatorPos := -1
	escaped := false

	lineLen := len(line)
	for curr := 0; curr < lineLen; curr++ {
		if !escaped {
			if line[curr] == '=' || line[curr] == ':' {
				separatorPos = curr
				break
			}

			// use first found space character as separator candidate if there is not one found yet
			// continue iteration as there might be another valid separator character further in the line
			sepFound := separatorPos > -1
			if !sepFound && isSpace(rune(line[curr])) {
				separatorPos = curr
			}
		}
		escaped = line[curr] == '\\'
	}
	return separatorPos
}

func parseValue(value string) (any, error) {
	if val, err := functions.DecodeStringValue(value); err == nil {
		return val, nil
	}
	return parseString(value)
}

func parseString(value string) (string, error) {
	builder := bytes.Buffer{}
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
					builder.WriteRune(decoded)
					offset += 4
					continue
				} else {
					builder.WriteRune(escapeCharacter(char))
				}
			} else {
				// Trailing backslash - treat as literal
				builder.WriteRune('\\')
			}
			offset++
		} else {
			decoded, decodedLen := utf8.DecodeRuneInString(value[offset:])
			if decoded == utf8.RuneError && decodedLen == 1 {
				return "", fmt.Errorf("invalid UTF-8 sequence at position %d", offset)
			}
			builder.WriteRune(decoded)
			offset += decodedLen
		}
	}
	return builder.String(), nil
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
	default:
		return rune(char)
	}
}

func extractKeys(line string, endPos int) ([]string, error) {
	keys := strings.TrimSpace(line[:endPos])
	parsedKeys, err := parseString(keys)
	if err != nil {
		return nil, fmt.Errorf("error parsing key '%s': %w", keys, err)
	}
	if parsedKeys == "" {
		return nil, fmt.Errorf("empty key found")
	}
	return strings.Split(parsedKeys, "."), nil
}

func extractValue(line string, startPos int) string {
	return strings.TrimSpace(line[startPos:])
}

func isEmpty(line string) bool {
	return len(line) == 0
}

func isNewline(line string) bool {
	return ((line[0] == '\n') || (line[0] == '\r'))
}

func isComment(line string) bool {
	return ((line[0] == '#') || (line[0] == '!'))
}

func isSpace(char rune) bool {
	return char == ' ' || char == '\t' || char == '\f'
}

func isContinuousLine(line string) bool {
	backslashesCount := 0
	for i := len(line) - 1; i >= 0 && line[i] == '\\'; i-- {
		backslashesCount++
	}
	return backslashesCount%2 == 1
}

func (d *Decoder) Decode() (keys []string, value any, err error) {
	line, err := d.bufReader.ReadString('\n')
	if err != nil {
		if err == io.EOF {
			return nil, nil, io.EOF
		}
		return nil, nil, fmt.Errorf("error reading line: %w", err)
	}

	line = strings.TrimSpace(line)
	if isEmpty(line) || isComment(line) || isNewline(line) {
		return nil, nil, nil
	}

	var separatorPos = findSeparatorPosition(line)
	if separatorPos == -1 {
		// key with no value
		keys, err = extractKeys(line, len(line))
		if err != nil {
			return nil, nil, fmt.Errorf("error extracting key from line '%s': %w", line, err)
		}
		return keys, "", nil // Return empty string instead of nil for consistency
	}

	keys, err = extractKeys(line, separatorPos)
	if err != nil {
		return nil, nil, fmt.Errorf("error extracting key from line '%s': %w", line, err)
	}

	valuePart := extractValue(line, separatorPos+1)
	isContinuous := isContinuousLine(valuePart)

	if isContinuous {
		builder := strings.Builder{}

		// Parse the first part
		parsedValue, err := parseString(valuePart)
		if err != nil {
			return nil, nil, fmt.Errorf("error parsing value '%s': %w", valuePart, err)
		}
		builder.WriteString(parsedValue)

		for isContinuous {
			line, err = d.bufReader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					// Handle end of file during continuation
					break
				}
				return nil, nil, fmt.Errorf("error reading continuation line: %w", err)
			}
			line = strings.TrimSpace(line)
			isContinuous = isContinuousLine(line)

			parsedLine, err := parseString(line)
			if err != nil {
				return nil, nil, fmt.Errorf("error parsing continuation line '%s': %w", line, err)
			}
			builder.WriteString(parsedLine)
		}
		value = builder.String()
	} else {
		value, err = parseValue(valuePart)
		if err != nil {
			return nil, nil, fmt.Errorf("error parsing value '%s': %w", valuePart, err)
		}
	}
	return keys, value, nil
}

func NewDecoder(bufReader *bufio.Reader) *Decoder {
	return &Decoder{bufReader: bufReader}
}
