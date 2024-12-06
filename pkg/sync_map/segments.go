package syncMap

import (
	"parallel-course-work/pkg/hash"
	linkedList "parallel-course-work/pkg/linked_list"
	"sync"
	"sync/atomic"
	"unsafe"
)

type Bucket[V any] struct {
	Key   string
	Value V
}

type Segment[V any] struct {
	innerArray    []*linkedList.LinkedList[Bucket[V]]
	maxLoadFactor float64
	size          atomic.Int64
	lock          sync.RWMutex
	isResizing    bool
	resizeCond    *sync.Cond
}

func NewSegment[V any](initialCapacity int, loadFactor float64) *Segment[V] {
	syncMap := &Segment[V]{
		innerArray:    make([]*linkedList.LinkedList[Bucket[V]], initialCapacity),
		maxLoadFactor: loadFactor,
		lock:          sync.RWMutex{},
		isResizing:    false,
	}
	syncMap.resizeCond = sync.NewCond(&syncMap.lock)
	return syncMap
}

func (h *Segment[V]) resizeSafe() error {
	h.lock.Lock()
	defer h.lock.Unlock()

	h.isResizing = true

	innerArrayCopy := make([]*linkedList.LinkedList[Bucket[V]], len(h.innerArray))
	copy(innerArrayCopy, h.innerArray)

	newSize := 2 * len(h.innerArray)
	h.innerArray = make([]*linkedList.LinkedList[Bucket[V]], newSize)
	h.size.Store(0)

	for _, list := range innerArrayCopy {
		for list != nil && list.GetSize() != 0 {
			element := list.RemoveFront()
			hashCode, err := hash.Calculate(element.Key)
			if err != nil {
				return err
			}

			idx := hashCode % uint64(len(h.innerArray))
			if h.innerArray[idx] == nil {
				h.innerArray[idx] = linkedList.NewWithInitValue[Bucket[V]](element)
			} else {
				h.innerArray[idx].AddFront(element)
			}

			h.size.Add(1)
		}
	}

	h.isResizing = false
	h.resizeCond.Broadcast()

	return nil
}

func (h *Segment[V]) needToResizeSafe() bool {
	h.lock.RLock()
	defer h.lock.RUnlock()
	loadFactor := float64(h.GetSize()) / float64(len(h.innerArray))
	return loadFactor > h.maxLoadFactor
}

func (h *Segment[V]) checkLoadFactorAndResizeSafe() error {
	if h.needToResizeSafe() {
		return h.resizeSafe()
	}

	return nil
}

func (h *Segment[V]) getBucketIndexFromKeySafe(key string) (uint64, error) {
	hashCode, err := hash.Calculate(key)
	if err != nil {
		return 0, err
	}

	h.lock.RLock()
	defer h.lock.RUnlock()
	idx := hashCode % uint64(len(h.innerArray))
	return idx, nil
}

func (h *Segment[V]) innerArrayIdxIsNilSafe(idx uint64) bool {
	h.lock.RLock()
	defer h.lock.RUnlock()
	return h.innerArray[idx] == nil
}

// NOT WORKING
func (h *Segment[V]) PutSafe(key string, value V) error {
	h.lock.Lock()
	for h.isResizing {
		h.resizeCond.Wait()
	}
	h.lock.Unlock()

	index, err := h.getBucketIndexFromKeySafe(key)
	if err != nil {
		return err
	}

	//First try to put the value in the empty bucket
	//If not successful, then try to put the value in the existing bucket
	if h.innerArrayIdxIsNilSafe(index) {
		newList := linkedList.NewWithInitValue[Bucket[V]](
			&Bucket[V]{
				Key:   key,
				Value: value,
			},
		)

		swapped := atomic.CompareAndSwapPointer(
			(*unsafe.Pointer)(unsafe.Pointer(&h.innerArray[index])),
			nil,
			unsafe.Pointer(newList))

		if swapped {
			h.size.Add(1)
			return nil
		}
	}

	//If the bucket is not empty, then try to find the key in the bucket
	{
		h.lock.Lock()

		for h.isResizing {
			h.resizeCond.Wait()
		}

		element, found := h.innerArray[index].FindByStructField("Key", key)
		if found {
			element.Value = value
			return nil
		}

		item := &Bucket[V]{key, value}
		h.innerArray[index].AddFront(item)
		h.size.Add(1)

		h.lock.Unlock()
	}

	return h.checkLoadFactorAndResizeSafe()
}

func (h *Segment[V]) GetSafe(key string) (*Bucket[V], bool) {
	index, err := h.getBucketIndexFromKeySafe(key)
	if err != nil {
		return nil, false
	}

	h.lock.RLock()
	defer h.lock.RUnlock()
	if h.innerArray[index] == nil {
		return nil, false
	}

	element, found := h.innerArray[index].FindByStructField("Key", key)
	if !found {
		return nil, false
	}

	return element, true
}

func (h *Segment[V]) RemoveSafe(key string) error {
	index, err := h.getBucketIndexFromKeySafe(key)
	if err != nil {
		return err
	}

	if h.innerArrayIdxIsNilSafe(index) {
		return nil
	}

	h.lock.Lock()
	defer h.lock.Unlock()
	err = h.innerArray[index].RemoveByStructField("Key", key)
	if err != nil {
		return err
	}

	h.size.Add(-1)
	return nil
}

func (h *Segment[V]) SetMaxLoadFactor(maxLoadFactor float64) {
	h.maxLoadFactor = maxLoadFactor
}

func (h *Segment[V]) GetLoadFactorSafe() float64 {
	h.lock.Lock()
	defer h.lock.Unlock()
	return float64(h.GetSize()) / float64(len(h.innerArray))
}

func (h *Segment[V]) GetSize() int64 {
	return h.size.Load()
}

func (h *Segment[V]) GetInnerArrayLenSafe() int {
	h.lock.RLock()
	defer h.lock.RUnlock()
	return len(h.innerArray)
}
