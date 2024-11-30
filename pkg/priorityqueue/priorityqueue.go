package priorityqueue

import (
	"sync"
)

type PriorityQueue[T any] struct {
	sync.RWMutex
	heap *heap[T]
}

func New[T any](compare CompareFunc[T]) *PriorityQueue[T] {
	return &PriorityQueue[T]{
		sync.RWMutex{},
		newHeap[T](compare),
	}
}

func (pq *PriorityQueue[T]) GetItems() []*T {
	pq.RLock()
	defer pq.RUnlock()
	return pq.heap.GetData()
}

func (pq *PriorityQueue[T]) Size() int {
	pq.RLock()
	defer pq.RUnlock()
	return pq.heap.Size()
}

func (pq *PriorityQueue[T]) Empty() bool {
	pq.RLock()
	defer pq.RUnlock()
	return pq.heap.Empty()
}

func (pq *PriorityQueue[T]) Push(element *T) {
	pq.Lock()
	defer pq.Unlock()
	pq.heap.Push(element)
}

func (pq *PriorityQueue[T]) Pop() (*T, error) {
	pq.Lock()
	defer pq.Unlock()
	return pq.heap.Pop()
}
