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
)

type YamlUnmarshal struct{}

func (yamlUnmarshal *YamlUnmarshal) Unmarshal(source serialization.SourceDataReader) (*ds.Node[string, interface{}], error) {
	var fileContent, err = source.Read()
	if err != nil {
		return nil, err
	}
	var content = yaml.MapSlice{}
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

func unmarshalItem(item yaml.MapItem, node *ds.Node[string, interface{}]) {
	if childItems, ok := item.Value.(yaml.MapSlice); ok {
		for _, childItem := range childItems {
			childItemKey := childItem.Key.(string)
			var childNodeRef = node.AddChild(childItemKey)
			unmarshalItem(childItem, childNodeRef)
		}
	} else {
		functions.UnmarshalValue(node, item.Value)
	}
}
