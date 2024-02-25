// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package execution

import (
	"fmt"
	"sync"

	"github.com/conformize/conformize/common"
	"github.com/conformize/conformize/common/ds"
	"github.com/conformize/conformize/common/path"
	"github.com/spaolacci/murmur3"
)

type valueRefHash [16]byte

func computePathStepsHashes(steps path.Steps) []valueRefHash {
	stepsLen := len(steps)
	hashes := make([]valueRefHash, stepsLen)

	var hash valueRefHash
	hasher := murmur3.New128()

	stepIdx := 0
	step := steps[stepIdx]
	hasher.Write([]byte(step.String()))
	hasher.Sum(hash[:0])
	hashes[stepIdx] = hash

	for stepIdx < stepsLen-1 {
		stepIdx++
		hasher.Write(hash[:])
		step = steps[stepIdx]
		if _, ok := step.(path.AttributeStep); ok {
			hasher.Write([]byte(":attribute:"))
		}

		hasher.Write([]byte(step.String()))
		hasher.Sum(hash[:0])
		hashes[stepIdx] = hash
	}
	return hashes
}

func computeRefHash(refAlias string) valueRefHash {
	var hash valueRefHash
	hasher := murmur3.New128()
	hasher.Write([]byte(refAlias))
	hasher.Sum(hash[:0])
	return hash
}

type ValueReferencesStore struct {
	valueReferences    map[valueRefHash]*ds.Node[string, interface{}]
	valuePathEvaluator common.ValuePathEvaluator
	rwLock             sync.RWMutex
}

func NewValueReferencesStore() *ValueReferencesStore {
	return &ValueReferencesStore{
		valueReferences: make(map[valueRefHash]*ds.Node[string, interface{}]),
	}
}

func (valRefStore *ValueReferencesStore) AddReference(refAlias string, valueRef *ds.Node[string, interface{}]) error {
	refHash := computeRefHash(refAlias)
	valRefStore.rwLock.Lock()
	defer valRefStore.rwLock.Unlock()
	if _, found := valRefStore.valueReferences[refHash]; found {
		return fmt.Errorf("value reference '%s' already defined", refAlias)
	}
	valRefStore.valueReferences[refHash] = valueRef
	return nil
}

func (valRefStore *ValueReferencesStore) GetReference(refAlias string) (*ds.Node[string, interface{}], bool) {
	refHash := computeRefHash(refAlias)
	valRefStore.rwLock.RLock()
	defer valRefStore.rwLock.RUnlock()
	refVal, found := valRefStore.valueReferences[refHash]
	return refVal, found
}

func (valRefStore *ValueReferencesStore) GetAtPath(refPath *path.Path) (*ds.Node[string, interface{}], error) {
	steps := refPath.Steps()

	stepHashes := computePathStepsHashes(steps)
	stepIdx := 0
	stepHash := stepHashes[stepIdx]

	valRefStore.rwLock.RLock()
	valueRef, found := valRefStore.valueReferences[stepHash]
	valRefStore.rwLock.RUnlock()
	if !found {
		root := steps[stepIdx]
		return nil, fmt.Errorf("couldn't find reference '%s'", root.String())
	}

	missingHashes := make(map[valueRefHash]*ds.Node[string, interface{}])
	for idx := len(stepHashes) - 1; idx > stepIdx; idx-- {
		stepHash = stepHashes[idx]
		valRefStore.rwLock.RLock()
		ref, found := valRefStore.valueReferences[stepHash]
		valRefStore.rwLock.RUnlock()

		if found {
			valueRef = ref
			stepIdx = idx + 1
			steps = steps[stepIdx:]
			break
		}
	}

	if fastForward := stepIdx > 0; !fastForward {
		stepIdx++
		steps.Next()
	}

	var walkErr error
	valueRef, walkErr = valRefStore.valuePathEvaluator.Evaluate(valueRef, path.NewPath(steps), func(step string, node *ds.Node[string, interface{}]) {
		stepHash = stepHashes[stepIdx]
		stepIdx++
		valRefStore.rwLock.RLock()
		if _, exists := valRefStore.valueReferences[stepHash]; !exists {
			missingHashes[stepHash] = node
		}
		valRefStore.rwLock.RUnlock()
	})

	valRefStore.rwLock.Lock()
	defer valRefStore.rwLock.Unlock()
	for hash, node := range missingHashes {
		valRefStore.valueReferences[hash] = node
	}

	return valueRef, walkErr
}
