// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package ds

import "sort"

const initialCapacity = 5

type DependencyGraph[T comparable] struct {
	edges   map[T][]T
	weights map[T]int

	index   int
	indices map[T]int
	lowLink map[T]int
	onStack map[T]bool
	stack   []T

	cycles [][]T
	order  []T
}

func NewDependencyGraph[T comparable]() *DependencyGraph[T] {
	return &DependencyGraph[T]{
		edges:   make(map[T][]T, initialCapacity),
		weights: make(map[T]int, initialCapacity),
		indices: make(map[T]int, initialCapacity),
		lowLink: make(map[T]int, initialCapacity),
		onStack: make(map[T]bool, initialCapacity),
	}
}

func (dg *DependencyGraph[T]) AddEdge(from, to T, weight int) {
	if _, exists := dg.weights[from]; !exists {
		dg.weights[from] = 0
	}

	if _, exists := dg.weights[to]; !exists {
		dg.weights[to] = weight
	}

	dg.edges[from] = append(dg.edges[from], to)
}

func (dg *DependencyGraph[T]) Run() {
	dg.index = 0
	dg.cycles = nil
	dg.order = make([]T, 0, len(dg.weights))
	dg.stack = make([]T, 0, len(dg.weights))

	for node := range dg.weights {
		if _, visited := dg.indices[node]; !visited {
			dg.strongConnect(node)
		}
	}

	if !dg.HasCycles() && len(dg.order) > 0 {
		levels := dg.computeLevels()

		sort.SliceStable(dg.order, func(i, j int) bool {
			nodeI, nodeJ := dg.order[i], dg.order[j]
			if levels[nodeI] != levels[nodeJ] {
				return levels[nodeI] < levels[nodeJ]
			}
			return dg.weights[nodeI] > dg.weights[nodeJ]
		})
	}
}

func (dg *DependencyGraph[T]) computeLevels() map[T]int {
	levels := make(map[T]int)

	for _, node := range dg.order {
		maxDepLevel := -1
		for _, dep := range dg.edges[node] {
			if levels[dep] > maxDepLevel {
				maxDepLevel = levels[dep]
			}
		}
		levels[node] = maxDepLevel + 1
	}

	return levels
}

func (dg *DependencyGraph[T]) strongConnect(v T) {
	dg.indices[v] = dg.index
	dg.lowLink[v] = dg.index
	dg.index++
	dg.stack = append(dg.stack, v)
	dg.onStack[v] = true

	for _, w := range dg.edges[v] {
		if _, ok := dg.indices[w]; !ok {
			dg.strongConnect(w)
			if dg.lowLink[w] < dg.lowLink[v] {
				dg.lowLink[v] = dg.lowLink[w]
			}
		} else if dg.onStack[w] {
			if dg.indices[w] < dg.lowLink[v] {
				dg.lowLink[v] = dg.indices[w]
			}
		}
	}

	if dg.lowLink[v] == dg.indices[v] {
		var scc []T
		for {
			w := dg.stack[len(dg.stack)-1]
			dg.stack = dg.stack[:len(dg.stack)-1]
			dg.onStack[w] = false
			scc = append(scc, w)
			if w == v {
				break
			}
		}

		if len(scc) > 1 {
			dg.cycles = append(dg.cycles, scc)
		} else {
			dg.order = append(dg.order, scc[0])
		}
	}
}

func (dg *DependencyGraph[T]) HasCycles() bool {
	return len(dg.cycles) > 0
}

func (dg *DependencyGraph[T]) GetCycles() [][]T {
	return dg.cycles
}

func (dg *DependencyGraph[T]) GetOrder() []T {
	return dg.order
}
