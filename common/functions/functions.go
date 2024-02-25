// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package functions

import "github.com/conformize/conformize/common/ds"

func MergeTrees[K comparable, V any](trees ...*ds.Node[K, V]) *ds.Node[K, V] {
	mergedTree := ds.NewNode[K, V]()
	for _, tree := range trees {
		for nodeKey, node := range tree.Children() {
			if existingNode, found := mergedTree.GetChild(nodeKey); !found {
				newNodeRef := mergedTree.AddChild(nodeKey)
				*newNodeRef = *ds.MergeNodes(newNodeRef, node)
			} else {
				*existingNode = *ds.MergeNodes(existingNode, node)
			}
		}
	}
	return mergedTree
}
