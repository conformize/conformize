// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package yaml

import (
	"fmt"

	"github.com/conformize/conformize/common/ds"
	"github.com/conformize/conformize/serialization"
	"github.com/conformize/conformize/serialization/unmarshal/functions"

	"gopkg.in/yaml.v2"
	yamlv2 "gopkg.in/yaml.v2"
)

type YamlUnmarshal struct {
	source serialization.SourceDataReader
}

func (yamlUnmarshal *YamlUnmarshal) Unmarshal() (*ds.Node[string, interface{}], error) {
	var fileContent, err = yamlUnmarshal.source.Read()
	if err != nil {
		return nil, err
	}
	var content = yamlv2.MapSlice{}
	err = yaml.Unmarshal(fileContent, &content)
	if err == nil {
		var rootNode = ds.NewNode[string, interface{}]()
		for _, item := range content {
			nodeKey := item.Key.(string)
			nodeRef := rootNode.AddChild(nodeKey)
			unmarshalItem(item, nodeRef)
		}
		return rootNode, nil
	}
	return nil, fmt.Errorf("failed to unmashal yaml content: %v", err.Error())
}

func unmarshalItem(item yamlv2.MapItem, node *ds.Node[string, interface{}]) {
	if childItems, ok := item.Value.(yamlv2.MapSlice); ok {
		for _, childItem := range childItems {
			childItemKey := childItem.Key.(string)
			var childNodeRef = node.AddChild(childItemKey)
			unmarshalItem(childItem, childNodeRef)
		}
	} else {
		functions.UnmarshalValue(node, item.Value)
	}
}

func NewYamlUnmarshal(source serialization.SourceDataReader) *YamlUnmarshal {
	return &YamlUnmarshal{source}
}