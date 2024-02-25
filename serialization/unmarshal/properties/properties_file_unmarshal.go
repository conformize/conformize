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
	"strings"

	"github.com/conformize/conformize/common/ds"
	"github.com/conformize/conformize/serialization"
)

type PropertiesFileUnmarshal struct{}

func (propUnmarshal *PropertiesFileUnmarshal) Unmarshal(source serialization.SourceDataReader) (*ds.Node[string, any], error) {
	fileContent, err := source.Read()
	if err != nil {
		return nil, err
	}

	bReader := bytes.NewReader(fileContent)
	rootNode := ds.NewNode[string, any]()
	decoder := NewDecoder(bufio.NewReader(bReader))

	lineNumber := 0
	for {
		keys, value, err := decoder.Decode()
		lineNumber++

		if err != nil {
			if err == io.EOF {
				break // Normal end of file
			}
			return nil, fmt.Errorf("error parsing properties file at line %d: %w", lineNumber, err)
		}

		if keys == nil {
			continue // Empty line or comment
		}

		nodeRef := rootNode
		keysLen := len(keys)
		for i := 0; i < keysLen; i++ {
			nodeKey := keys[i]
			if nodeKey == "" {
				return nil, fmt.Errorf("empty key component found in key path '%s' at line %d", strings.Join(keys, "."), lineNumber)
			}

			nodes, ok := nodeRef.GetChildren(nodeKey)
			if ok {
				nodeRef = nodes.First()
				continue
			}
			nodeRef = nodeRef.AddChild(nodeKey)
		}
		nodeRef.Value = value
	}
	return rootNode, nil
}
