package ds

import "container/heap"

type QueueItem[T any] struct {
	Value    T
	Priority int
}

type priorityComparatorFunc[T any] func(a, b *QueueItem[T]) bool

func MinPriorityComparator[T any](a, b *QueueItem[T]) bool {
	return a.Priority < b.Priority
}

func MaxPriorityComparator[T any](a, b *QueueItem[T]) bool {
	return a.Priority > b.Priority
}

type priorityQueue[T any] struct {
	items              []*QueueItem[T]
	priorityComparator priorityComparatorFunc[T]
}

func (pq *priorityQueue[T]) Len() int {
	return len(pq.items)
}

func (pq *priorityQueue[T]) Less(i, j int) bool {
	return pq.priorityComparator(pq.items[i], pq.items[j])
}

func (pq *priorityQueue[T]) Swap(i, j int) {
	pq.items[i], pq.items[j] = pq.items[j], pq.items[i]
}

func (pq *priorityQueue[T]) Push(item any) {
	qItem := item.(*QueueItem[T])
	pq.items = append(pq.items, qItem)
}

func (pq *priorityQueue[T]) Pop() any {
	old := pq.items
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	pq.items = old[:n-1]
	return item
}

func (pq *priorityQueue[T]) IsEmpty() bool {
	return len(pq.items) == 0
}

func newPriorityQueue[T any](comparator priorityComparatorFunc[T]) *priorityQueue[T] {
	var cmp priorityComparatorFunc[T] = MinPriorityComparator[T]
	if comparator != nil {
		cmp = comparator
	}

	pq := &priorityQueue[T]{
		items:              []*QueueItem[T]{},
		priorityComparator: cmp,
	}
	return pq
}

type PriorityQueue[T any] struct {
	queue *priorityQueue[T]
}

func (pq *PriorityQueue[T]) Push(priority int, item T) {
	heap.Push(pq.queue, &QueueItem[T]{Value: item, Priority: priority})
}

func (pq *PriorityQueue[T]) Peek() (T, bool) {
	if pq.IsEmpty() {
		var zero T
		return zero, false
	}
	item := pq.queue.items[0]
	return item.Value, true
}

func (pq *PriorityQueue[T]) IsEmpty() bool {
	return pq.queue.IsEmpty()
}

func (pq *PriorityQueue[T]) Pop() (T, bool) {
	if pq.IsEmpty() {
		var zero T
		return zero, false
	}
	qItem := heap.Pop(pq.queue).(*QueueItem[T])
	return qItem.Value, true
}

func (pq *PriorityQueue[T]) Len() int {
	return pq.queue.Len()
}

func NewPriorityQueue[T any](comparator priorityComparatorFunc[T]) *PriorityQueue[T] {
	pq := &PriorityQueue[T]{
		queue: newPriorityQueue(comparator),
	}
	heap.Init(pq.queue)
	return pq
}
