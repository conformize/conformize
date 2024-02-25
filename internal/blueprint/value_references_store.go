// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package blueprint

import (
	"fmt"
	"sync"

	"github.com/conformize/conformize/common"
	"github.com/conformize/conformize/common/path"
	"github.com/conformize/conformize/common/pathparser"
	"github.com/conformize/conformize/internal/blueprint/elements"
)

type ValueReferencesStore struct {
	valueReferences map[string]*elements.ValueReference
	pathParser      *pathparser.PathParser
	valueWalker     *common.ValuePathWalker
	rwLock          *sync.RWMutex
}

func NewValueReferencesStore() *ValueReferencesStore {
	return &ValueReferencesStore{
		valueReferences: make(map[string]*elements.ValueReference),
		pathParser:      pathparser.NewPathParser(),
		valueWalker:     &common.ValuePathWalker{},
		rwLock:          &sync.RWMutex{},
	}
}

func (valRefStore *ValueReferencesStore) AddReference(refAlias string, valueRef *elements.ValueReference) error {
	valRefStore.rwLock.Lock()
	defer valRefStore.rwLock.Unlock()
	if _, found := valRefStore.valueReferences[refAlias]; found {
		return fmt.Errorf("value reference '%s' already exists", refAlias)
	}
	valRefStore.valueReferences[refAlias] = valueRef
	return nil
}

func (valRefStore *ValueReferencesStore) GetReference(refAlias string) (*elements.ValueReference, bool) {
	valRefStore.rwLock.RLock()
	defer valRefStore.rwLock.RUnlock()
	refVal, found := valRefStore.valueReferences[refAlias]
	return refVal, found
}

func (valRefStore *ValueReferencesStore) GetAtPath(pathStr string) (*elements.ValueReference, error) {
	steps, err := valRefStore.pathParser.Parse(pathStr)
	if err != nil {
		return nil, err
	}
	root := steps[0]
	if rootRef, found := valRefStore.GetReference(root.String()); !found {
		return nil, fmt.Errorf("value reference %s at path %s not found", root.String(), root.String())
	} else {
		pathToWalk := path.NewPath(steps[1:])
		value, walkErr := valRefStore.valueWalker.Walk(rootRef.Node, pathToWalk)
		if walkErr != nil {
			return nil, walkErr
		}
		return &elements.ValueReference{
			Node: value,
			Meta: &elements.ValueReferenceMеta{
				ValuePath: pathStr,
			},
		}, nil
	}
}
