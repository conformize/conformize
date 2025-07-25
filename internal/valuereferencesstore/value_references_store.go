// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package valuereferencesstore

import (
	"fmt"
	"sync"

	"github.com/conformize/conformize/common"
	"github.com/conformize/conformize/common/ds"
	"github.com/conformize/conformize/common/path"
)

type ValueReferencesStore struct {
	valueReferences   map[string]*ds.Node[string, interface{}]
	exprPathEvaluator common.ExpressionPathEvaluator
	rwLock            sync.RWMutex
}

var (
	instance *ValueReferencesStore
	once     sync.Once
)

func newValueReferencesStore() *ValueReferencesStore {
	return &ValueReferencesStore{
		valueReferences: make(map[string]*ds.Node[string, interface{}]),
	}
}

func Instance() *ValueReferencesStore {
	once.Do(func() {
		instance = newValueReferencesStore()
	})
	return instance
}

func (valRefStore *ValueReferencesStore) AddReference(refAlias string, valueRef *ds.Node[string, interface{}]) error {
	valRefStore.rwLock.Lock()
	defer valRefStore.rwLock.Unlock()
	if _, found := valRefStore.valueReferences[refAlias]; found {
		return fmt.Errorf("value reference '%s' already defined", refAlias)
	}
	valRefStore.valueReferences[refAlias] = valueRef
	return nil
}

func (valRefStore *ValueReferencesStore) GetReference(refAlias string) (*ds.Node[string, interface{}], bool) {
	valRefStore.rwLock.RLock()
	defer valRefStore.rwLock.RUnlock()
	refVal, found := valRefStore.valueReferences[refAlias]
	return refVal, found
}

func (valRefStore *ValueReferencesStore) GetAtPath(refPath *path.Path) (*ds.Node[string, interface{}], error) {
	steps := refPath.Steps()
	if len(steps) == 0 {
		return nil, fmt.Errorf("path cannot be empty")
	}

	root, _ := steps.Next()

	valRefStore.rwLock.RLock()
	valueRef, found := valRefStore.valueReferences[root.String()]
	valRefStore.rwLock.RUnlock()
	if !found {
		return nil, fmt.Errorf("couldn't find reference '%s'", root.String())
	}

	var walkErr error
	valueRef, walkErr = valRefStore.exprPathEvaluator.Evaluate(valueRef, path.NewPath(steps))
	return valueRef, walkErr
}
