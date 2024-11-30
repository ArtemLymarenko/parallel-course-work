package priorityqueue

import (
	"errors"
)

type CompareFunc[T any] func(a *T, b *T) bool

type heap[T any] struct {
	data    []*T
	compare CompareFunc[T]
}

func newHeap[T any](compare CompareFunc[T]) *heap[T] {
	return &heap[T]{make([]*T, 0), compare}
}

func (h *heap[T]) swap(firstIdx int, secondIdx int) {
	temp := h.data[firstIdx]
	h.data[firstIdx] = h.data[secondIdx]
	h.data[secondIdx] = temp
}

func (h *heap[T]) GetData() []*T {
	return h.data
}

func (h *heap[T]) Size() int {
	return len(h.data)
}

func (h *heap[T]) Empty() bool {
	return len(h.data) == 0
}

func (h *heap[T]) Push(element *T) {
	h.data = append(h.data, element)
	lastIdx := h.Size() - 1
	h.siftUp(lastIdx)
}

func (h *heap[T]) Pop() (removed *T, err error) {
	heapSize := h.Size()
	if heapSize == 0 {
		return removed, errors.New("heap is empty")
	}

	lastIdx := heapSize - 1
	h.swap(0, lastIdx)

	removed = h.data[lastIdx]
	h.data = h.data[:lastIdx]
	h.siftDown(0)

	return removed, nil
}

func (h *heap[T]) siftDown(idx int) {
	for {
		leftChildIdx := 2*idx + 1
		rightChildIdx := 2*idx + 2

		smallest := idx
		if leftChildIdx < h.Size() && h.compare(h.data[smallest], h.data[leftChildIdx]) {
			smallest = leftChildIdx
		}

		if rightChildIdx < h.Size() && h.compare(h.data[smallest], h.data[rightChildIdx]) {
			smallest = rightChildIdx
		}

		if smallest == idx {
			break
		}

		h.swap(idx, smallest)
		idx = smallest
	}
}

func (h *heap[T]) siftUp(idx int) {
	for idx > 0 {
		parentIdx := (idx - 1) / 2
		if h.compare(h.data[idx], h.data[parentIdx]) {
			break
		}

		h.swap(idx, parentIdx)
		idx = parentIdx
	}
}
