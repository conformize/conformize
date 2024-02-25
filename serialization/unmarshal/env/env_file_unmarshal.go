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

	"github.com/conformize/conformize/common/ds"
	"github.com/conformize/conformize/serialization"
)

type EnvFileUnmarshal struct {
	fileSource *serialization.FileSource
}

func (envUnmarshal *EnvFileUnmarshal) Unmarshal() (*ds.Node[string, interface{}], error) {
	var fileContent, err = envUnmarshal.fileSource.Read()
	if err == nil {
		bReader := bytes.NewReader(fileContent)
		rootNode := ds.NewNode[string, interface{}]()
		decoder := NewDecoder(bufio.NewReader(bReader))
		for key, value, err := decoder.Decode(); err == nil; key, value, err = decoder.Decode() {
			if key == nil {
				continue
			}
			nodeRef := rootNode.AddChild(*key)
			nodeRef.Value = value
		}
		return rootNode, nil
	}
	return nil, err
}

func NewEnvFileUnmarshal(fileSource *serialization.FileSource) *EnvFileUnmarshal {
	return &EnvFileUnmarshal{fileSource: fileSource}
}
