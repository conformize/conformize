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
	"strings"
	"unicode/utf16"

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
	for pos := 0; pos < len(*key); pos++ {
		char := (*key)[pos]
		if !(isLowercase(char) || isUppercase(char) ||
			((pos > 0) && (isDigit(char) || isUnderscore(char)))) {
			return false
		}
	}
	return true
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

func parseValue(value string) interface{} {
	if val, err := functions.DecodeStringValue(value); err == nil {
		return val
	}
	return parseString(value)
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
	case 'b':
		return '\b'
	default:
		return rune(char)
	}
}

func parseString(value string) string {
	buf := bytes.Buffer{}
	valLen := len(value)
	offset := 0

	var char byte
	var nextPos = 0
	for offset < valLen {
		char = value[offset]
		if char == '\\' {
			nextPos = offset + 1
			if nextPos < valLen {
				offset = nextPos
				char = value[offset]
				if char == 'u' {
					offset++
					hexPart := value[offset : offset+4]
					decoded := decodeUnicode(hexPart)
					buf.Write([]byte(string(decoded)))
					offset += 4
					continue
				} else {
					buf.Write([]byte(string(escapeCharacter(char))))
					offset++
					continue
				}
			}
		} else if char == '#' {
			break
		}
		offset++
		buf.WriteByte(char)
	}
	return buf.String()
}

func (d *Decoder) Decode() (key *string, value interface{}, err error) {
	line, err := d.bufReader.ReadString('\n')
	if err != nil {
		return nil, nil, errEOF
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
		return nil, nil, err
	}

	separatorPos++
	valuePart := extractValue(lineTrim, separatorPos)
	if isMultilineToken(valuePart) {
		token := &valuePart
		closed := false

		bldr := strings.Builder{}
		for valuePart, err = d.bufReader.ReadString('\n'); !closed && err == nil; valuePart, err = d.bufReader.ReadString('\n') {
			valuePart = strings.TrimSpace(valuePart)
			if isMultilineToken(valuePart) {
				closed = *token == valuePart
				continue
			}
			bldr.WriteByte('\n')
			valuePart = parseString(strings.TrimSpace(valuePart))
			bldr.WriteString(valuePart)
		}

		if !closed {
			return nil, nil, fmt.Errorf("malformed multiline value! Expected token %s not found", *token)
		}
		return key, bldr.String(), nil
	}
	return key, parseValue(valuePart), nil
}

func NewDecoder(bufReader *bufio.Reader) *Decoder {
	return &Decoder{bufReader: bufReader}
}
