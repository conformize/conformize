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

func parseValue(value string) interface{} {
	if val, err := functions.DecodeStringValue(value); err == nil {
		return val
	}
	return parseString(value)
}

func parseString(value string) string {
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
					var hexPart = value[offset : offset+4]
					var decoded = decodeUnicode(hexPart)
					builder.WriteRune(decoded)
					offset += 4
					continue
				} else {
					builder.WriteRune(escapeCharacter(char))
				}
			}
			offset++
		} else {
			decoded, decodedLen := utf8.DecodeRuneInString(value[offset:])
			builder.WriteRune(decoded)
			offset += decodedLen
		}
	}
	return builder.String()
}

func parseHexString(hexStr string) uint16 {
	var res uint16
	for _, r := range hexStr {
		res *= 16
		if '0' <= r && r <= '9' {
			res += uint16(r - '0')
		} else if 'a' <= r && r <= 'f' {
			res += uint16(r - 'a' + 10)
		} else if 'A' <= r && r <= 'F' {
			res += uint16(r - 'A' + 10)
		}
	}
	return res
}

func decodeUnicode(val string) rune {
	return utf16.Decode([]uint16{parseHexString(val[0:4])})[0]
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

func extractKeys(line string, endPos int) []string {
	keys := strings.TrimSpace(line[:endPos])
	keys = parseString(keys)
	return strings.Split(keys, ".")
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

func (d *Decoder) Decode() (keys []string, value interface{}, err error) {
	line, err := d.bufReader.ReadString('\n')
	if err != nil {
		return nil, nil, fmt.Errorf("EOF")
	}

	line = strings.TrimSpace(line)
	if isEmpty(line) || isComment(line) || isNewline(line) {
		return nil, nil, nil
	}

	var separatorPos = findSeparatorPosition(line)
	if separatorPos == -1 {
		// key with no value
		keys = extractKeys(line, len(line))
		return keys, nil, nil
	}
	keys = extractKeys(line, separatorPos)
	valuePart := extractValue(line, separatorPos+1)
	isContinuous := isContinuousLine(valuePart)

	if isContinuous {
		builder := strings.Builder{}
		builder.WriteString(parseString(valuePart))

		for isContinuous {
			line, _ = d.bufReader.ReadString('\n')
			line = strings.TrimSpace(line)
			isContinuous = isContinuousLine(line)
			builder.WriteString(parseString(line))
		}
		value = builder.String()
	} else {
		value = parseValue(valuePart)
	}
	return keys, value, nil
}

func NewDecoder(bufReader *bufio.Reader) *Decoder {
	return &Decoder{bufReader: bufReader}
}
