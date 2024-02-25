// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package toml

import (
	"bytes"

	"github.com/conformize/conformize/common/ds"
	"github.com/conformize/conformize/serialization"
	"github.com/conformize/conformize/serialization/unmarshal/functions"

	"github.com/BurntSushi/toml"
)

type TomlFilelUnmarshal struct {
	fileSource *serialization.FileSource
}

func (tomlUnmarshal *TomlFilelUnmarshal) Unmarshal() (*ds.Node[string, interface{}], error) {
	fileContent, err := tomlUnmarshal.fileSource.Read()
	if err != nil {
		return nil, err
	}

	data := map[string]interface{}{}
	_, err = toml.NewDecoder(bytes.NewReader(fileContent)).Decode(&data)
	if err != nil {
		return nil, err
	}

	rootNode := ds.NewNode[string, interface{}]()
	for key, value := range data {
		nodeRef := rootNode.AddChild(key)
		functions.UnmarshalValue(nodeRef, value)
	}
	return rootNode, nil

}

func NewTomlFileUnmarshal(fileSource *serialization.FileSource) *TomlFilelUnmarshal {
	return &TomlFilelUnmarshal{fileSource: fileSource}
}
