package syncMap

import (
	"hash/maphash"
	"parallel-course-work/pkg/hash"
	linkedList "parallel-course-work/pkg/linked_list"
	"sync"
)

type Bucket[V any] struct {
	Key   string
	Value *V
}

const maxLoadFactor float64 = 0.9

type segment[V any] struct {
	innerArray []*linkedList.LinkedList[Bucket[V]]
	lock       sync.RWMutex
	size       int64
	seed       maphash.Seed
}

func NewSegment[V any](initialCapacity int) *segment[V] {
	seg := &segment[V]{
		innerArray: make([]*linkedList.LinkedList[Bucket[V]], initialCapacity),
		lock:       sync.RWMutex{},
		seed:       hash.MakeRandomSeed(),
	}

	return seg
}

func (h *segment[V]) resize() {
	h.lock.Lock()
	h.seed = hash.MakeRandomSeed()
	newCapacity := 2 * len(h.innerArray)
	resizeArray := make([]*linkedList.LinkedList[Bucket[V]], newCapacity)

	for index := range h.innerArray {
		for h.innerArray[index] != nil && h.innerArray[index].GetSize() > 0 {
			item := h.innerArray[index].RemoveFront()

			newIdx := int64(hash.Get(h.seed, item.Key) % uint64(len(resizeArray)))
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

func (h *segment[V]) PutSafe(bucket *Bucket[V]) bool {
	h.lock.Lock()
	exists := h.putBucket(bucket)
	if !exists {
		h.size++
	}
	h.lock.Unlock()

	h.checkAndStartResize()

	return exists
}

func (h *segment[V]) putBucket(bucket *Bucket[V]) (exists bool) {
	hashCode := hash.Get(h.seed, bucket.Key)

	index := int64(hashCode % uint64(len(h.innerArray)))
	if h.innerArray[index] == nil {
		h.innerArray[index] = linkedList.NewWithInitValue[Bucket[V]](bucket)
		return false
	}

	element, found := h.innerArray[index].Find(func(current *Bucket[V]) bool {
		return current.Key == bucket.Key
	})
	if found {
		element.Value = bucket.Value
		return true
	}

	h.innerArray[index].AddFront(bucket)
	return false
}

func (h *segment[V]) GetSafe(key string) (*Bucket[V], bool) {
	h.lock.RLock()
	defer h.lock.RUnlock()

	index := int64(hash.Get(h.seed, key) % uint64(len(h.innerArray)))
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

// ModifySafe changes existing value of map using callback function
// but if after the modification value is nil then the entire bucket will be deleted
func (h *segment[V]) ModifySafe(key string, cb func(modify *V) *V) (modified bool, removed bool) {
	h.lock.Lock()
	defer h.lock.Unlock()

	index := int64(hash.Get(h.seed, key) % uint64(len(h.innerArray)))
	if h.innerArray[index] == nil {
		return false, false
	}

	element, found := h.innerArray[index].Find(func(current *Bucket[V]) bool {
		return current.Key == key
	})
	if found {
		element.Value = cb(element.Value)
		if element.Value == nil {
			h.innerArray[index].Remove(func(current *Bucket[V]) bool {
				return current.Key == key
			})
			h.size--
			return true, true
		}
		return true, false
	} else {
		return false, false
	}
}

func (h *segment[V]) RemoveSafe(key string) {
	h.lock.Lock()
	defer h.lock.Unlock()

	index := int64(hash.Get(h.seed, key) % uint64(len(h.innerArray)))
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
