package syncMap

import (
	"parallel-course-work/pkg/hash"
	linkedList "parallel-course-work/pkg/linked_list"
	"sync"
)

const (
	BucketFieldKey = "Key"
)

type Bucket struct {
	Key   string
	Value []string
}

type linkedListArray []linkedList.ILinkedList[Bucket]

type syncHashMap struct {
	innerArray    linkedListArray
	locks         []sync.RWMutex
	maxLoadFactor float64
	size          int
	lockSize      int
}

func NewSyncHashMap(size int, loadFactor float64, lockSize int) *syncHashMap {
	syncMap := &syncHashMap{
		innerArray:    make(linkedListArray, size),
		maxLoadFactor: loadFactor,
	}

	for i := 0; i < lockSize; i++ {
		syncMap.locks = append(syncMap.locks, sync.RWMutex{})
	}

	return syncMap
}

func (h *syncHashMap) SetMaxLoadFactor(maxLoadFactor float64) {
	h.maxLoadFactor = maxLoadFactor
}

func (h *syncHashMap) GetLoadFactor() float64 {
	return float64(h.size) / float64(len(h.innerArray))
}

func (h *syncHashMap) GetSize() int {
	return h.size
}

func (h *syncHashMap) getBucketIndexFromKey(key string) (uint64, error) {
	hashCode, err := hash.Calculate(key)
	if err != nil {
		return 0, err
	}

	return hashCode % uint64(len(h.innerArray)), nil
}

func (h *syncHashMap) resizeMap() error {
	innerArrayCopy := make(linkedListArray, len(h.innerArray))
	copy(innerArrayCopy, h.innerArray)

	newSize := int((h.maxLoadFactor * 2) * float64(len(h.innerArray)))
	h.innerArray = make(linkedListArray, newSize)
	h.size = 0

	for _, list := range innerArrayCopy {
		for list != nil && list.GetSize() != 0 {
			element := list.RemoveFront()
			err := h.Insert(element.Key, element.Value)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (h *syncHashMap) checkLoadFactorAndResize() error {
	if h.GetLoadFactor() > h.maxLoadFactor {
		return h.resizeMap()
	}

	return nil
}

func (h *syncHashMap) Insert(key string, value []string) error {
	index, err := h.getBucketIndexFromKey(key)
	if err != nil {
		return err
	}

	if h.innerArray[index] == nil {
		h.innerArray[index] = linkedList.New[Bucket]()
	}

	list := h.innerArray[index]
	element, found := list.FindByStructField(BucketFieldKey, key)
	if found {
		element.Value = value
		return nil
	}

	list.AddFront(&Bucket{key, value})
	h.size++

	err = h.checkLoadFactorAndResize()
	return err
}

func (h *syncHashMap) Get(key string) (*Bucket, error) {
	index, err := h.getBucketIndexFromKey(key)
	if err != nil {
		return nil, err
	}

	if h.innerArray[index] == nil {
		return nil, ErrElementNotFound
	}

	element, found := h.innerArray[index].FindByStructField(BucketFieldKey, key)
	if !found {
		return nil, ErrElementNotFound
	}

	return element, nil
}

func (h *syncHashMap) Remove(key string) error {
	index, err := h.getBucketIndexFromKey(key)
	if err != nil {
		return err
	}

	if h.innerArray[index] == nil {
		return nil
	}

	err = h.innerArray[index].RemoveByStructField(BucketFieldKey, key)
	if err != nil {
		return err
	}

	h.size--
	return nil
}
