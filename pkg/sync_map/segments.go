package syncMap

import (
	"parallel-course-work/pkg/hash"
	linkedList "parallel-course-work/pkg/linked_list"
	"sync"
	"sync/atomic"
)

type Bucket[V any] struct {
	Key   string
	Value V
}

const maxSegments = 16
const maxLoadFactor float64 = 0.75

type syncMap[V any] struct {
	innerArray     []*linkedList.LinkedList[Bucket[V]]
	resizeArray    []*linkedList.LinkedList[Bucket[V]]
	locks          []sync.RWMutex
	resizeLock     sync.RWMutex
	isResizing     bool
	size           atomic.Int64
	resizeCtlIndex atomic.Int64
	resizeRoutines atomic.Int32
	resizeCond     *sync.Cond
}

func New[V any](initialCapacity int) *syncMap[V] {
	locks := make([]sync.RWMutex, maxSegments)
	sMap := &syncMap[V]{
		innerArray: make([]*linkedList.LinkedList[Bucket[V]], initialCapacity),
		locks:      locks,
		resizeLock: sync.RWMutex{},
		isResizing: false,
	}
	sMap.size.Store(0)
	sMap.resizeCtlIndex.Store(0)
	sMap.resizeRoutines.Store(0)
	sMap.resizeCond = sync.NewCond(&sMap.resizeLock)
	return sMap
}

func (h *syncMap[V]) GetLoadFactor() float64 {
	h.resizeLock.RLock()
	defer h.resizeLock.RUnlock()
	return float64(h.size.Load()) / float64(len(h.innerArray))
}

func (h *syncMap[V]) getBucketIndexFromKeySafe(key string) int64 {
	hashCode, err := hash.Calculate(key)
	if err != nil {
		return 0
	}

	h.resizeLock.RLock()
	defer h.resizeLock.RUnlock()
	return int64(hashCode % uint64(len(h.innerArray)))
}

func (h *syncMap[V]) getResizeBucketIndexFromKeySafe(key string) int64 {
	hashCode, err := hash.Calculate(key)
	if err != nil {
		return 0
	}

	h.resizeLock.RLock()
	defer h.resizeLock.RUnlock()
	if len(h.resizeArray) == 0 {
		return -1
	}
	return int64(hashCode % uint64(len(h.resizeArray)))
}

func (h *syncMap[V]) GetInnerArrayLenSafe() int64 {
	h.resizeLock.RLock()
	defer h.resizeLock.RUnlock()
	return int64(len(h.innerArray))
}

func (h *syncMap[V]) getResizeArrayLenSafe() int64 {
	h.resizeLock.RLock()
	defer h.resizeLock.RUnlock()
	return int64(len(h.resizeArray))
}

func (h *syncMap[V]) resize() {
	h.resizeRoutines.Add(1)
	for {
		index := h.resizeCtlIndex.Add(1) - 1
		if index >= h.GetInnerArrayLenSafe() {
			h.resizeRoutines.Add(-1)
			break
		}

		lockIdx := int(index % maxSegments)
		h.locks[lockIdx].Lock()
		for h.innerArray[index] != nil && h.innerArray[index].GetSize() > 0 {
			item := h.innerArray[index].RemoveFront()
			h.locks[lockIdx].Unlock()

			newIdx := h.getResizeBucketIndexFromKeySafe(item.Key)
			newLockIdx := newIdx % maxSegments

			h.locks[newLockIdx].Lock()
			if h.resizeArray[newIdx] == nil {
				h.resizeArray[newIdx] = linkedList.NewWithInitValue[Bucket[V]](item)
			} else {
				h.resizeArray[newIdx].AddFront(item)
			}
			h.locks[newLockIdx].Unlock()

			h.locks[lockIdx].Lock()
		}

		h.locks[lockIdx].Unlock()
	}

	h.resizeLock.Lock()
	defer h.resizeLock.Unlock()
	if h.isResizing && h.resizeRoutines.Load() == 0 {
		h.innerArray = h.resizeArray
		h.resizeArray = nil
		h.resizeCtlIndex.Store(0)
		h.isResizing = false
		h.resizeCond.Broadcast()
	}
}

func (h *syncMap[V]) checkAndStartResize() {
	h.resizeLock.Lock()
	defer h.resizeLock.Unlock()

	loadFactor := float64(h.size.Load()) / float64(len(h.innerArray))
	if loadFactor < maxLoadFactor || h.isResizing {
		return
	}

	newResizeArrayLen := 2 * len(h.innerArray)
	h.resizeArray = make([]*linkedList.LinkedList[Bucket[V]], newResizeArrayLen)
	h.isResizing = true

	go h.resize()
}

func (h *syncMap[V]) waitIfResizing() {
	h.resizeLock.Lock()
	for h.isResizing {
		h.resizeCond.Wait()
	}
	h.resizeLock.Unlock()
}

func (h *syncMap[V]) checkIsResizing() bool {
	h.resizeLock.Lock()
	defer h.resizeLock.Unlock()
	return h.isResizing
}

func (h *syncMap[V]) PutSafe(key string, value V) {
	h.waitIfResizing()

	if !h.checkIsResizing() {
		index := h.getBucketIndexFromKeySafe(key)
		h.putToBucket(h.innerArray, index, key, value)
	} else {
		index := h.getResizeBucketIndexFromKeySafe(key)
		if index == -1 {
			index = h.getBucketIndexFromKeySafe(key)
			h.putToBucket(h.innerArray, index, key, value)
		} else {
			h.putToBucket(h.resizeArray, index, key, value)
		}
	}

	h.checkAndStartResize()
}

func (h *syncMap[V]) putToBucket(
	array []*linkedList.LinkedList[Bucket[V]],
	index int64,
	key string,
	value V,
) {
	lockIdx := int(index % maxSegments)
	h.locks[lockIdx].Lock()
	defer func() {
		h.locks[lockIdx].Unlock()
		h.size.Add(1)
	}()

	if array[index] == nil {
		array[index] = linkedList.NewWithInitValue[Bucket[V]](&Bucket[V]{key, value})
		return
	}

	element, found := array[index].FindByStructField("Key", key)
	if found {
		element.Value = value
	} else {
		array[index].AddFront(&Bucket[V]{Key: key, Value: value})
	}
}

func (h *syncMap[V]) GetSafe(key string) (*Bucket[V], bool) {
	index := h.getBucketIndexFromKeySafe(key)

	lockIdx := int(index % maxSegments)
	h.locks[lockIdx].RLock()
	defer h.locks[lockIdx].RUnlock()

	if h.innerArray[index] == nil {
		return nil, false
	}

	element, found := h.innerArray[index].FindByStructField("Key", key)
	if !found {
		return nil, false
	}

	return element, true
}

func (h *syncMap[V]) RemoveSafe(key string) error {
	index := h.getBucketIndexFromKeySafe(key)

	lockIdx := int(index % maxSegments)
	h.locks[lockIdx].Lock()
	defer h.locks[lockIdx].Unlock()

	if h.innerArray[index] == nil {
		return nil
	}

	const fieldName = "Key"
	err := h.innerArray[index].RemoveByStructField(fieldName, key)
	if err != nil {
		return err
	}

	h.size.Add(-1)
	return nil
}

func (h *syncMap[V]) GetLoadFactorSafe() float64 {
	h.resizeLock.Lock()
	defer h.resizeLock.Unlock()
	return float64(h.GetSize()) / float64(len(h.innerArray))
}

func (h *syncMap[V]) GetSize() int64 {
	return h.size.Load()
}
