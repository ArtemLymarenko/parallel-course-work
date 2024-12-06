package syncMap

//
//import (
//	"parallel-course-work/pkg/hash"
//	linkedList "parallel-course-work/pkg/linked_list"
//	"sync"
//)
//
//type Bucket[V any] struct {
//	Key   string
//	Value V
//}
//
//type linkedListArray[V any] []linkedList.ILinkedList[Bucket[V]]
//
//type syncHashMap[V any] struct {
//	innerArray    linkedListArray[V]
//	locks         []sync.RWMutex
//	maxLoadFactor float64
//	size          int
//	lockSize      int
//}
//
//func NewSyncHashMap[V any](initialCapacity int, loadFactor float64) *syncHashMap[V] {
//	defaultLockSize := 16
//	return &syncHashMap[V]{
//		innerArray:    make(linkedListArray[V], initialCapacity),
//		locks:         make([]sync.RWMutex, defaultLockSize),
//		maxLoadFactor: loadFactor,
//	}
//}
//
//func (h *syncHashMap[V]) SetMaxLoadFactor(maxLoadFactor float64) {
//	h.maxLoadFactor = maxLoadFactor
//}
//
//func (h *syncHashMap[V]) GetLoadFactor() float64 {
//	return float64(h.size) / float64(len(h.innerArray))
//}
//
//func (h *syncHashMap[V]) GetSize() int {
//	return h.size
//}
//
//func (h *syncHashMap[V]) getBucketIndexFromKey(key string) (uint64, error) {
//	hashCode, err := hash.Calculate(key)
//	if err != nil {
//		return 0, err
//	}
//
//	return hashCode % uint64(len(h.innerArray)), nil
//}
//
//func (h *syncHashMap[V]) resizeMap() error {
//	innerArrayCopy := make(linkedListArray[V], len(h.innerArray))
//	copy(innerArrayCopy, h.innerArray)
//
//	newSize := int((h.maxLoadFactor * 2) * float64(len(h.innerArray)))
//	h.innerArray = make(linkedListArray[V], newSize)
//	h.size = 0
//
//	for _, list := range innerArrayCopy {
//		for list != nil && list.GetSize() != 0 {
//			element := list.RemoveFront()
//			err := h.Put(element.Key, element.Value)
//			if err != nil {
//				return err
//			}
//		}
//	}
//
//	return nil
//}
//
//func (h *syncHashMap[V]) checkLoadFactorAndResize() error {
//	if h.GetLoadFactor() > h.maxLoadFactor {
//		return h.resizeMap()
//	}
//
//	return nil
//}
//
//func (h *syncHashMap[V]) Put(key string, value V) error {
//	index, err := h.getBucketIndexFromKey(key)
//	if err != nil {
//		return err
//	}
//
//	if h.innerArray[index] == nil {
//		h.innerArray[index] = linkedList.New[Bucket[V]]()
//		item := &Bucket[V]{key, value}
//		h.innerArray[index].AddFront(item)
//		h.size++
//		return h.checkLoadFactorAndResize()
//	}
//
//	element, found := h.innerArray[index].FindByStructField("Key", key)
//	if found {
//		element.Value = value
//	}
//
//	return h.checkLoadFactorAndResize()
//}
//
//func (h *syncHashMap[V]) Get(key string) (*Bucket[V], bool) {
//	index, err := h.getBucketIndexFromKey(key)
//	if err != nil {
//		return nil, false
//	}
//
//	if h.innerArray[index] == nil {
//		return nil, false
//	}
//
//	element, found := h.innerArray[index].FindByStructField("Key", key)
//	if !found {
//		return nil, false
//	}
//
//	return element, true
//}
//
//func (h *syncHashMap[V]) Remove(key string) error {
//	index, err := h.getBucketIndexFromKey(key)
//	if err != nil {
//		return err
//	}
//
//	if h.innerArray[index] == nil {
//		return nil
//	}
//
//	err = h.innerArray[index].RemoveByStructField("Key", key)
//	if err != nil {
//		return err
//	}
//
//	h.size--
//	return nil
//}
