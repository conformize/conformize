// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package ds

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v2"
)

type NodeList[K comparable, V any] []*Node[K, V]

func (nl NodeList[K, V]) Count() int {
	return len(nl)
}

func (nl NodeList[K, V]) First() *Node[K, V] {
	if len(nl) == 0 {
		return nil
	}
	return nl[0]
}

func (nl NodeList[K, V]) Get(index int) *Node[K, V] {
	if index < 0 || index >= len(nl) {
		return nil
	}
	return nl[index]
}

type Node[K comparable, V any] struct {
	Key        K
	Value      V
	attributes map[string]*Node[string, any]
	Parent     *Node[K, V]
	children   map[K]NodeList[K, V]
}

func NewNode[K comparable, V any]() *Node[K, V] {
	var val V
	return &Node[K, V]{
		Value:      val,
		attributes: make(map[string]*Node[string, any]),
		children:   make(map[K]NodeList[K, V]),
	}
}

func (n *Node[K, V]) AddChild(key K) *Node[K, V] {
	child := NewNode[K, V]()
	child.Key = key
	child.Parent = n
	n.children[key] = append(n.children[key], child)
	return child
}

func (n *Node[K, V]) Append(key K, child *Node[K, V]) {
	if child == nil {
		return
	}
	if n.children == nil {
		n.children = make(map[K]NodeList[K, V])
	}
	child.Key = key
	child.Parent = n
	n.children[key] = append(n.children[key], child)
}

func (n *Node[K, V]) GetParent() *Node[K, V] {
	return n.Parent
}

func (n *Node[K, V]) GetChildren(key K) (NodeList[K, V], bool) {
	children, found := n.children[key]
	return children, found
}

func (n *Node[K, V]) Children() map[K]NodeList[K, V] {
	if n.children != nil {
		return n.children
	}
	return map[K]NodeList[K, V]{}
}

func (n *Node[K, V]) AddAttribute(attributeName string, attributeValue any) {
	n.attributes[attributeName] = NewNode[string, any]()
	n.attributes[attributeName].Value = attributeValue
}

func (n *Node[K, V]) GetAttribute(attributeName string) (*Node[string, any], bool) {
	attribute, found := n.attributes[attributeName]
	return attribute, found
}

func (n *Node[K, V]) GetAttributes() map[string]*Node[string, any] {
	if n.attributes != nil {
		return n.attributes
	}
	return map[string]*Node[string, any]{}
}

func MergeNodes[K comparable, V any](node *Node[K, V], oNode *Node[K, V]) *Node[K, V] {
	newNode := NewNode[K, V]()
	if node == nil && oNode == nil {
		return newNode
	}

	newNode.Key = node.Key
	newNode.Value = oNode.Value

	if node.Parent != nil {
		newNode.Parent = node.Parent
	}

	mergeAttributes(newNode, node)
	mergeAttributes(newNode, oNode)

	mergedChildren := mergeChildren(node, oNode)
	for childKey, childList := range mergedChildren {
		for _, child := range childList {
			newNode.Append(childKey, child)
		}
	}

	return newNode
}

func mergeAttributes[K comparable, V any](node *Node[K, V], oNode *Node[K, V]) {
	for attrKey, attr := range oNode.attributes {
		node.AddAttribute(attrKey, attr.Value)
	}
}

func mergeChildren[K comparable, V any](node *Node[K, V], oNode *Node[K, V]) map[K]NodeList[K, V] {
	merged := make(map[K]NodeList[K, V])

	for childKey, childList := range node.Children() {
		for _, child := range childList {
			merged[childKey] = append(merged[childKey], child.Clone())
		}
	}

	for childKey, oChildList := range oNode.Children() {
		if existingList, found := merged[childKey]; found {
			for i, oChild := range oChildList {
				if i < len(existingList) {
					merged[childKey][i] = MergeNodes(existingList[i], oChild)
				} else {
					merged[childKey] = append(merged[childKey], oChild.Clone())
				}
			}
		} else {
			for _, oChild := range oChildList {
				merged[childKey] = append(merged[childKey], oChild.Clone())
			}
		}
	}
	return merged
}

func (n *Node[K, V]) PrintTree() {
	for nodeKey, nodeList := range n.children {
		for _, node := range nodeList {
			fmt.Printf("Node: %v, Value: %v\n", nodeKey, node.Value)
			printAttributes(node, 1)
			printNodes(node, 1)
		}
	}
}

func (n *Node[K, V]) Clone() *Node[K, V] {
	cloned := NewNode[K, V]()
	cloned.Key = n.Key
	cloned.Value = n.Value

	for k, attr := range n.GetAttributes() {
		cloned.AddAttribute(k, attr.Value)
	}

	for k, children := range n.Children() {
		for _, child := range children {
			cloned.Append(k, child.Clone())
		}
	}

	return cloned
}

func printNodes[K comparable, V any](n *Node[K, V], level int) {
	for childKey, children := range n.Children() {
		for _, child := range children {
			fmt.Printf("%sNode: %v, Value: %v\n", strings.Repeat(" ", level*2), childKey, child.Value)
			printAttributes(child, level+1)
			printNodes(child, level+1)
		}
	}
}

func printAttributes[K comparable, V any](n *Node[K, V], level int) {
	for key, value := range n.GetAttributes() {
		fmt.Printf("%sAttribute: %s, Value: %v\n", strings.Repeat(" ", level*2), key, value.Value)
	}
}

func (n *Node[K, V]) UnmarshalYAML(unmarshal func(any) error) error {
	var content yaml.MapSlice
	if err := unmarshal(&content); err != nil {
		return fmt.Errorf("failed to unmarshal yaml content: %v", err)
	}

	for _, item := range content {
		key := item.Key.(K)
		child := n.AddChild(key)
		unmarshalYAMLItem(item.Value, child)
	}
	return nil
}

func unmarshalYAMLItem[K comparable, V any](value any, node *Node[K, V]) {
	switch val := value.(type) {
	case yaml.MapSlice:
		for _, item := range val {
			key := item.Key.(K)
			child := node.AddChild(key)
			unmarshalYAMLItem(item.Value, child)
		}
	case []any:
		for _, elem := range val {
			if innerMap, ok := elem.(yaml.MapSlice); ok {
				for _, innerItem := range innerMap {
					key := innerItem.Key.(K)
					child := node.AddChild(key)
					unmarshalYAMLItem(innerItem.Value, child)
				}
			} else {
				unmarshalValue(node, val)
				break
			}
		}
	default:
		unmarshalValue(node, val)
	}
}

func unmarshalValue[K comparable, V any](nodeRef *Node[K, V], value any) {
	if valMap, ok := value.(map[K]any); ok {
		for key, v := range valMap {
			childNode := nodeRef.AddChild(key)
			unmarshalValue(childNode, v)
		}
	} else {
		nodeRef.Value = value.(V)
	}
}
