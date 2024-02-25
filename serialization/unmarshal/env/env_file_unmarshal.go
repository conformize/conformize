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
	"fmt"
	"io"

	"github.com/conformize/conformize/common/ds"
	"github.com/conformize/conformize/serialization"
)

type EnvFileUnmarshal struct{}

func (envUnmarshal *EnvFileUnmarshal) Unmarshal(source serialization.SourceDataReader) (*ds.Node[string, any], error) {
	fileContent, err := source.Read()
	if err != nil {
		return nil, err
	}

	bReader := bytes.NewReader(fileContent)
	rootNode := ds.NewNode[string, any]()
	decoder := NewDecoder(bufio.NewReader(bReader))

	lineNumber := 0
	for {
		key, value, err := decoder.Decode()
		lineNumber++

		if err != nil {
			if err == io.EOF {
				break // Normal end of file
			}
			return nil, fmt.Errorf("error parsing .env file at line %d: %w", lineNumber, err)
		}

		if key == nil {
			continue // Empty line or comment
		}

		if *key == "" {
			return nil, fmt.Errorf("empty environment variable name at line %d", lineNumber)
		}

		nodeRef := rootNode.AddChild(*key)
		nodeRef.Value = value
	}
	return rootNode, nil
}
