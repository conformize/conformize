// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package xml

import (
	"bytes"
	"encoding/xml"
	"io"

	"github.com/conformize/conformize/common/ds"
	"github.com/conformize/conformize/serialization"
	"github.com/conformize/conformize/serialization/unmarshal/functions"
)

type XmlFileUnmarshal struct{}

func (xmlUnmarshal *XmlFileUnmarshal) Unmarshal(source serialization.SourceDataReader) (*ds.Node[string, any], error) {
	var fileContent, err = source.Read()
	if err != nil {
		return nil, err
	}
	decoder := xml.NewDecoder(io.NopCloser(bytes.NewReader(fileContent)))
	rootNode := ds.NewNode[string, any]()
	nodeRef := rootNode

	for token, err := decoder.Token(); err == nil; token, err = decoder.Token() {
		switch tokenType := token.(type) {
		case xml.StartElement:
			var nodeKey = tokenType.Name.Local
			if nodeRef != nil {
				nodeRef = nodeRef.AddChild(nodeKey)
			} else {
				nodeRef = rootNode.AddChild(nodeKey)
			}
			for _, attr := range tokenType.Attr {
				nodeRef.AddAttribute(attr.Name.Local, attr.Value)
			}
		case xml.EndElement:
			nodeRef = nodeRef.GetParent()
		case xml.CharData:
			if nodeRef != nil {
				var strValue = string(tokenType)
				if functions.IsWhiteSpace(strValue) {
					nodeRef.Value = nil
					continue
				}

				decodedValue, err := functions.DecodeStringValue(strValue)
				if err == nil {
					nodeRef.Value = decodedValue
					continue
				}
				nodeRef.Value = strValue
			}
		}
	}
	return rootNode, nil
}
