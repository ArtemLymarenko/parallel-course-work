package syncMap

//
//import (
//	"parallel-course-work/pkg/hash"
//	linkedList "parallel-course-work/pkg/linked_list"
//	"sync"
//	"sync/atomic"
//)
//
//type Bucket[V any] struct {
//	Key   string
//	Value V
//}
//
//const maxSegments = 16
//const maxLoadFactor float64 = 0.75
//
//type syncMap[V any] struct {
//	innerArray     []*linkedList.LinkedList[Bucket[V]]
//	resizeArray    []*linkedList.LinkedList[Bucket[V]]
//	locks          []sync.RWMutex
//	resizeLocks    []sync.RWMutex
//	lock           sync.RWMutex
//	resizeLock     sync.RWMutex
//	isResizing     bool
//	size           atomic.Int64
//	resizeCtlIndex atomic.Int64
//	resizeRoutines atomic.Int32
//	resizeCond     *sync.Cond
//}
//
//func NewSyncHashMap[V any](initialCapacity int, segments int) *syncMap[V] {
//	locks := make([]sync.RWMutex, maxSegments)
//	resizeLocks := make([]sync.RWMutex, maxSegments)
//	sMap := &syncMap[V]{
//		innerArray:  make([]*linkedList.LinkedList[Bucket[V]], initialCapacity),
//		locks:       locks,
//		lock:        sync.RWMutex{},
//		resizeLock:  sync.RWMutex{},
//		resizeLocks: resizeLocks,
//		isResizing:  false,
//	}
//	sMap.size.Store(0)
//	sMap.resizeCtlIndex.Store(0)
//	sMap.resizeRoutines.Store(0)
//	sMap.resizeCond = sync.NewCond(&sMap.resizeLock)
//
//	return sMap
//}
//
//func (h *syncMap[V]) getBucketIndexFromKeySafe(key string) int64 {
//	hashCode := hash.Get(key)
//	h.resizeLock.RLock()
//	defer h.resizeLock.RUnlock()
//	return int64(hashCode % uint64(len(h.innerArray)))
//}
//
//func (h *syncMap[V]) resize() {
//	h.resizeLock.RLock()
//	h.lock.RLock()
//
//	h.resizeRoutines.Add(1)
//	for {
//		index := h.resizeCtlIndex.Add(1) - 1
//		if index >= int64(len(h.innerArray)) || h.resizeArray == nil {
//			h.resizeRoutines.Add(-1)
//			break
//		}
//
//		lockIdx := int(index % maxSegments)
//		h.locks[lockIdx].Lock()
//		for h.innerArray[index] != nil && h.innerArray[index].GetSize() > 0 {
//			item := h.innerArray[index].RemoveFront()
//
//			newIdx := int64(hash.Get(item.Key) % uint64(len(h.resizeArray)))
//			newLockIdx := int(newIdx % maxSegments)
//
//			h.resizeLocks[newLockIdx].Lock()
//			if h.resizeArray[newIdx] == nil {
//				h.resizeArray[newIdx] = linkedList.NewWithInitValue[Bucket[V]](item)
//			} else {
//				h.resizeArray[newIdx].AddFront(item)
//			}
//			h.resizeLocks[newLockIdx].Unlock()
//		}
//		h.locks[lockIdx].Unlock()
//	}
//
//	h.lock.RUnlock()
//	h.resizeLock.RUnlock()
//
//	h.resizeLock.Lock()
//	if h.isResizing && h.resizeRoutines.Load() == 0 {
//		h.lock.Lock()
//		h.innerArray = h.resizeArray
//		h.lock.Unlock()
//
//		h.resizeArray = nil
//		h.isResizing = false
//		h.resizeCtlIndex.Store(0)
//		h.resizeCond.Broadcast()
//	}
//	h.resizeLock.Unlock()
//}
//
//func (h *syncMap[V]) checkAndStartResize() {
//	h.resizeLock.Lock()
//	h.lock.Lock()
//	loadFactor := float64(h.size.Load()) / float64(len(h.innerArray))
//	if loadFactor < maxLoadFactor || h.isResizing {
//		h.lock.Unlock()
//		h.resizeLock.Unlock()
//		return
//	}
//
//	newResizeCap := 2 * len(h.innerArray)
//	h.lock.Unlock()
//
//	h.resizeArray = make([]*linkedList.LinkedList[Bucket[V]], newResizeCap)
//	h.isResizing = true
//
//	h.resizeLock.Unlock()
//
//	go h.resize()
//}
//
//func (h *syncMap[V]) Put(key string, value V) {
//	h.resizeLock.Lock()
//	if h.isResizing {
//		h.resizeLock.Unlock()
//		h.resize()
//		h.resizeLock.Lock()
//	}
//
//	for h.isResizing {
//		h.resizeCond.Wait()
//	}
//	h.resizeLock.Unlock()
//
//	h.putToBucket(key, value)
//
//	h.checkAndStartResize()
//}
//
//func (h *syncMap[V]) putToBucket(
//	key string,
//	value V,
//) {
//	h.lock.RLock()
//	defer h.lock.RUnlock()
//
//	index := hash.Get(key) % uint64(len(h.innerArray))
//	lockIdx := int(index % maxSegments)
//	h.locks[lockIdx].Lock()
//	defer h.locks[lockIdx].Unlock()
//
//	if h.innerArray[index] == nil {
//		h.innerArray[index] = linkedList.NewWithInitValue[Bucket[V]](&Bucket[V]{key, value})
//		h.size.Add(1)
//		return
//	}
//
//	element, found := h.innerArray[index].Find(func(current *Bucket[V]) bool {
//		return current.Key == key
//	})
//	if found {
//		element.Value = value
//	} else {
//		h.innerArray[index].AddFront(&Bucket[V]{Key: key, Value: value})
//		h.size.Add(1)
//	}
//}
//
//func (h *syncMap[V]) Get(key string) (*Bucket[V], bool) {
//	h.lock.RLock()
//	defer h.lock.RUnlock()
//
//	index := hash.Get(key) % uint64(len(h.innerArray))
//	lockIdx := int(index % maxSegments)
//	h.locks[lockIdx].Lock()
//	defer h.locks[lockIdx].Unlock()
//
//	if h.innerArray[index] == nil {
//		return nil, false
//	}
//
//	element, found := h.innerArray[index].Find(func(current *Bucket[V]) bool {
//		return current.Key == key
//	})
//	if !found {
//		return nil, false
//	}
//
//	return element, true
//}
//
//func (h *syncMap[V]) Remove(key string) {
//	index := h.getBucketIndexFromKeySafe(key)
//
//	lockIdx := int(index % maxSegments)
//	h.locks[lockIdx].Lock()
//	defer h.locks[lockIdx].Unlock()
//
//	if h.innerArray[index] == nil {
//		return
//	}
//
//	h.innerArray[index].Remove(func(current *Bucket[V]) bool {
//		return current.Key == key
//	})
//
//	h.size.Add(-1)
//}
//
//func (h *syncMap[V]) GetSize() int64 {
//	return h.size.Load()
//}
