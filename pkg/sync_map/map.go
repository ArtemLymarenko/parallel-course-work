package syncMap

import (
	"hash/fnv"
	linkedList "parallel-course-work/pkg/linked_list"
	"strings"
)

type node[V any] struct {
	key   string
	value *V
}

type bucket[V any] []linkedList.ILinkedList[node[V]]

type syncHashMap[V any] struct {
	innerArray    bucket[V]
	maxLoadFactor float64
	size          int
}

func NewSyncHashMap[V any](size int, loadFactor float64) *syncHashMap[V] {
	return &syncHashMap[V]{
		innerArray:    make(bucket[V], size),
		maxLoadFactor: loadFactor,
	}
}

func (h *syncHashMap[V]) SetMaxLoadFactor(maxLoadFactor float64) {
	h.maxLoadFactor = maxLoadFactor
}

func (h *syncHashMap[V]) hash(key string) (uint64, error) {
	if len(strings.Trim(key, " ")) == 0 {
		return 0, ErrCalculatingHash
	}

	totalHash := fnv.New64a()
	_, err := totalHash.Write([]byte(key))
	if err != nil {
		return 0, ErrCalculatingHash
	}

	return totalHash.Sum64(), nil
}

func (h *syncHashMap[V]) GetInnerArray() bucket[V] {
	return h.innerArray
}

func (h *syncHashMap[V]) GetLoadFactor() float64 {
	return float64(h.size) / float64(len(h.innerArray))
}

func (h *syncHashMap[V]) GetSize() int {
	return h.size
}

func (h *syncHashMap[V]) resizeMap() error {
	innerArrayCopy := make(bucket[V], len(h.innerArray))
	copy(innerArrayCopy, h.innerArray)

	newSize := int((h.maxLoadFactor * 2) * float64(len(h.innerArray)))
	h.innerArray = make(bucket[V], newSize)
	h.size = 0

	for _, list := range innerArrayCopy {
		for list != nil && list.GetSize() != 0 {
			element := list.RemoveByIndex(0)
			err := h.Insert(element.key, element.value)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (h *syncHashMap[V]) Insert(key string, value *V) error {
	hashCode, err := h.hash(key)
	if err != nil {
		return err
	}

	index := hashCode % uint64(len(h.innerArray))
	if h.innerArray[index] == nil {
		h.innerArray[index] = linkedList.New[node[V]]()
	}

	newNode := node[V]{key, nil}
	_, err = h.innerArray[index].Find(newNode, func(a, b node[V]) bool {
		return a.key == b.key
	})
	if err != nil {
		h.innerArray[index].AddFront(node[V]{key, value})
		h.size++
	}

	if h.GetLoadFactor() > h.maxLoadFactor {
		err = h.resizeMap()

		if err != nil {
			return err
		}
	}

	return nil
}

func (h *syncHashMap[V]) Get(key string) (*V, error) {
	hashCode, hashErr := h.hash(key)
	if hashErr != nil {
		return nil, hashErr
	}

	index := hashCode % uint64(len(h.innerArray))
	if h.innerArray[index] == nil {
		return nil, ErrElementNotFound
	}

	newNode := node[V]{key, nil}
	element, err := h.innerArray[index].Find(newNode, func(a, b node[V]) bool {
		return a.key == b.key
	})

	if err != nil {
		return nil, err
	}

	return element.value, nil
}

func (h *syncHashMap[V]) Remove(key string) error {
	hashCode, hashErr := h.hash(key)
	if hashErr != nil {
		return hashErr
	}

	index := hashCode % uint64(len(h.innerArray))
	if h.innerArray[index] == nil {
		return nil
	}

	newNode := node[V]{key, nil}
	err := h.innerArray[index].Remove(newNode, func(a, b node[V]) bool {
		return a.key == b.key
	})

	if err != nil {
		return err
	}

	h.size--
	return nil
}
