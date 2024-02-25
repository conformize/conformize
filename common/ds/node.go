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

type Node[K comparable, V any] struct {
	Key        K
	Value      V
	attributes map[string]*Node[string, interface{}]
	parent     *Node[K, V]
	next       *Node[K, V]
	children   map[K]*Node[K, V]
}

func NewNode[K comparable, V any]() *Node[K, V] {
	var key K
	var val V
	return &Node[K, V]{
		Key:        key,
		Value:      val,
		next:       nil,
		attributes: make(map[string]*Node[string, interface{}]),
		children:   make(map[K]*Node[K, V]),
	}
}

func (n *Node[K, V]) AddChild(key K) *Node[K, V] {
	newNode := NewNode[K, V]()
	newNode.Key = key
	newNode.parent = n
	if childNode, found := n.children[key]; !found {
		n.children[key] = newNode
	} else {
		for childNode.next != nil {
			childNode = childNode.next
		}
		childNode.next = newNode
	}
	return newNode
}

func (n *Node[K, V]) GetParent() *Node[K, V] {
	return n.parent
}

func (n *Node[K, V]) Next() *Node[K, V] {
	return n.next
}

func (n *Node[K, V]) GetChild(key K) (*Node[K, V], bool) {
	child, found := n.children[key]
	return child, found
}

func (n *Node[K, V]) Children() map[K]*Node[K, V] {
	if n.children != nil {
		return n.children
	}
	return map[K]*Node[K, V]{}
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
	newNode.Key = oNode.Key
	newNode.Value = oNode.Value
	if node.parent != nil {
		newNode.parent = node.parent
	}
	mergeAttributes(newNode, node)
	mergeAttributes(newNode, oNode)
	children := mergeChildren(node, oNode)
	for _, child := range children {
		newChildRef := newNode.AddChild(child.Key)
		*newChildRef = *MergeNodes(newChildRef, child)
	}
	return newNode
}

func mergeAttributes[K comparable, V any](node *Node[K, V], oNode *Node[K, V]) {
	for _, attr := range oNode.attributes {
		node.AddAttribute(attr.Key, attr.Value)
	}
}

func mergeChildren[K comparable, V any](node *Node[K, V], oNode *Node[K, V]) []*Node[K, V] {
	merged := make(map[K]*Node[K, V])
	children := make([]*Node[K, V], 0)

	for _, child := range node.Children() {
		merged[child.Key] = child.Clone()
		children = append(children, merged[child.Key])
	}

	for _, child := range oNode.Children() {
		if existingNode, found := merged[child.Key]; !found {
			merged[child.Key] = child.Clone()
			children = append(children, merged[child.Key])
		} else {
			*merged[child.Key] = *MergeNodes(existingNode, child)
		}
	}
	return children
}

func (n *Node[K, V]) PrintTree() {
	for _, node := range n.children {
		fmt.Printf("Node: %v, Value: %v\n", node.Key, node.Value)
		printAttributes(node, 1)
		printNodes(node, 1)
	}
}

func (n *Node[K, V]) Clone() *Node[K, V] {
	return MergeNodes(NewNode[K, V](), n)
}

func printNodes[K comparable, V any](n *Node[K, V], level int) {
	for _, child := range n.Children() {
		fmt.Printf("%sNode: %v, Value: %v\n", strings.Repeat(" ", level*2), child.Key, child.Value)
		printAttributes(child, level+1)
		printNodes(child, level+1)
		for ch := child.Next(); ch != nil; ch = ch.Next() {
			fmt.Printf("%sNode: %v, Value: %v\n", strings.Repeat(" ", level*2), child.Key, child.Value)
			printAttributes(ch, level+1)
			printNodes(ch, level+1)
		}
	}
}

func printAttributes[K comparable, V any](n *Node[K, V], level int) {
	for key, value := range n.GetAttributes() {
		fmt.Printf("%sAttribute: %s, Value: %v\n", strings.Repeat(" ", level*2), key, value.Value)
	}
}
