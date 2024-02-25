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
	attributes map[string]*Node[string, interface{}]
	Parent     *Node[K, V]
	children   map[K]NodeList[K, V]
}

func NewNode[K comparable, V any]() *Node[K, V] {
	var val V
	return &Node[K, V]{
		Value:      val,
		attributes: make(map[string]*Node[string, interface{}]),
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

func (n *Node[K, V]) AddAttribute(attributeName string, attributeValue interface{}) {
	n.attributes[attributeName] = NewNode[string, interface{}]()
	n.attributes[attributeName].Value = attributeValue
}

func (n *Node[K, V]) GetAttribute(attributeName string) (*Node[string, interface{}], bool) {
	attribute, found := n.attributes[attributeName]
	return attribute, found
}

func (n *Node[K, V]) GetAttributes() map[string]*Node[string, interface{}] {
	if n.attributes != nil {
		return n.attributes
	}
	return map[string]*Node[string, interface{}]{}
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
