// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package ds

type DependencyGraph[T comparable] struct {
	nodes   map[T][]T
	index   int
	stack   []T
	indices map[T]int
	lowLink map[T]int
	onStack map[T]struct{}
	cycles  [][]T
	order   []T
}

func NewDependencyGraph[T comparable]() *DependencyGraph[T] {
	return &DependencyGraph[T]{
		nodes:   make(map[T][]T),
		indices: make(map[T]int),
		lowLink: make(map[T]int),
		onStack: make(map[T]struct{}),
	}
}

func (dg *DependencyGraph[T]) AddEdge(v, w T) {
	dg.nodes[v] = append(dg.nodes[v], w)
}

func (dg *DependencyGraph[T]) strongConnect(v T) {
	dg.indices[v] = dg.index
	dg.lowLink[v] = dg.index
	dg.index++
	dg.stack = append(dg.stack, v)
	dg.onStack[v] = struct{}{}

	for _, w := range dg.nodes[v] {
		if _, ok := dg.indices[w]; !ok {
			dg.strongConnect(w)
			dg.lowLink[v] = min(dg.lowLink[v], dg.lowLink[w])
		} else if _, onStack := dg.onStack[w]; onStack {
			dg.lowLink[v] = min(dg.lowLink[v], dg.indices[w])
		}
	}

	if dg.lowLink[v] == dg.indices[v] {
		var scc []T
		for {
			w := dg.stack[len(dg.stack)-1]
			dg.stack = dg.stack[:len(dg.stack)-1]
			delete(dg.onStack, w)
			scc = append(scc, w)
			if w == v {
				break
			}
		}
		if len(scc) > 1 {
			dg.cycles = append(dg.cycles, scc)
		} else {
			dg.order = append(dg.order, scc...)
		}
	}
}

func (dg *DependencyGraph[T]) HasCycles() bool {
	return len(dg.cycles) > 0
}

func (dg *DependencyGraph[T]) Run() {
	dg.index = 0
	for node := range dg.nodes {
		if _, ok := dg.indices[node]; !ok {
			dg.strongConnect(node)
		}
	}
}

func (dg *DependencyGraph[T]) IsEmpty() bool {
	return len(dg.nodes) == 0
}

func (dg *DependencyGraph[T]) GetCycles() [][]T {
	return dg.cycles
}

func (dg *DependencyGraph[T]) GetOrder() []T {
	return dg.order
}
