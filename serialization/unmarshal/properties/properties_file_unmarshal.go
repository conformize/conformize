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

	"github.com/conformize/conformize/common/ds"
	"github.com/conformize/conformize/serialization"
)

type PropertiesFileUnmarshal struct{}

func (propUnmarshal *PropertiesFileUnmarshal) Unmarshal(source serialization.SourceDataReader) (*ds.Node[string, interface{}], error) {
	fileContent, err := source.Read()
	if err != nil {
		return nil, err
	}

	bReader := bytes.NewReader(fileContent)
	rootNode := ds.NewNode[string, interface{}]()
	decoder := NewDecoder(bufio.NewReader(bReader))
	for keys, value, err := decoder.Decode(); err == nil; keys, value, err = decoder.Decode() {
		if keys == nil {
			continue
		}

		nodeRef := rootNode
		keysLen := len(keys)
		for i := 0; i < keysLen; i++ {
			nodeKey := keys[i]
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
