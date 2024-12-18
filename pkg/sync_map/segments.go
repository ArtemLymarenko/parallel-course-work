package syncMap

import (
	"parallel-course-work/pkg/hash"
	linkedList "parallel-course-work/pkg/linked_list"
	"sync"
)

type Bucket[V any] struct {
	Key   string
	Value V
	hash  uint64
}

func (b *Bucket[V]) GetHash() (uint64, bool) {
	return b.hash, b.hash != 0
}

const maxLoadFactor float64 = 0.75

type segment[V any] struct {
	innerArray []*linkedList.LinkedList[Bucket[V]]
	lock       sync.RWMutex
	size       int64
}

func NewSegment[V any](initialCapacity int) *segment[V] {
	seg := &segment[V]{
		innerArray: make([]*linkedList.LinkedList[Bucket[V]], initialCapacity),
		lock:       sync.RWMutex{},
	}

	return seg
}

func (h *segment[V]) resize() {
	h.lock.Lock()
	newCapacity := 2 * len(h.innerArray)
	resizeArray := make([]*linkedList.LinkedList[Bucket[V]], newCapacity)

	for index := range h.innerArray {
		for h.innerArray[index] != nil && h.innerArray[index].GetSize() > 0 {
			item := h.innerArray[index].RemoveFront()
			hashCode, ok := item.GetHash()
			if !ok {
				hashCode, _ = hash.Calculate(item.Key)
			}
			newIdx := int64(hashCode % uint64(newCapacity))
			if resizeArray[newIdx] == nil {
				resizeArray[newIdx] = linkedList.NewWithInitValue[Bucket[V]](item)
			} else {
				resizeArray[newIdx].AddFront(item)
			}
		}
	}

	h.innerArray = resizeArray
	resizeArray = nil
	h.lock.Unlock()
}

func (h *segment[V]) checkAndStartResize() {
	h.lock.Lock()
	loadFactor := float64(h.size) / float64(len(h.innerArray))
	if loadFactor < maxLoadFactor {
		h.lock.Unlock()
		return
	}
	h.lock.Unlock()
	h.resize()
}

func (h *segment[V]) PutSafe(bucket *Bucket[V]) {
	h.lock.Lock()
	h.putBucket(bucket)
	h.size++
	h.lock.Unlock()

	h.checkAndStartResize()
}

func (h *segment[V]) putBucket(bucket *Bucket[V]) {
	hashCode, ok := bucket.GetHash()
	if !ok {
		hashCode, _ = hash.Calculate(bucket.Key)
	}

	index := int64(hashCode % uint64(len(h.innerArray)))
	if h.innerArray[index] == nil {
		h.innerArray[index] = linkedList.NewWithInitValue[Bucket[V]](bucket)
		return
	}

	element, found := h.innerArray[index].Find(func(current *Bucket[V]) bool {
		return current.Key == bucket.Key
	})

	if found {
		element.Value = bucket.Value
	} else {
		h.innerArray[index].AddFront(bucket)
	}
}

func (h *segment[V]) GetSafe(hash uint64, key string) (*Bucket[V], bool) {
	h.lock.RLock()
	defer h.lock.RUnlock()

	index := int64(hash % uint64(len(h.innerArray)))
	if h.innerArray[index] == nil {
		return nil, false
	}

	element, found := h.innerArray[index].Find(func(current *Bucket[V]) bool {
		return current.Key == key
	})
	if !found {
		return nil, false
	}

	return element, true
}

func (h *segment[V]) RemoveSafe(hash uint64, key string) {
	h.lock.Lock()
	defer h.lock.Unlock()

	index := int64(hash % uint64(len(h.innerArray)))
	if h.innerArray[index] == nil {
		return
	}

	h.innerArray[index].Remove(func(current *Bucket[V]) bool {
		return current.Key == key
	})
	h.size--
}

func (h *segment[V]) GetSize() int64 {
	h.lock.RLock()
	defer h.lock.RUnlock()
	return h.size
}
