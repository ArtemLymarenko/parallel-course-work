package syncMap

import (
	"github.com/ArtemLymarenko/parallel-course-work/pkg/hash"
	linkedList "github.com/ArtemLymarenko/parallel-course-work/pkg/linked_list"
	"github.com/ArtemLymarenko/parallel-course-work/pkg/set"
	"hash/maphash"
	"sync"
)

type Bucket struct {
	Key   string
	Value *set.Set[string]
}

const maxLoadFactor float64 = 0.75

type segment struct {
	innerArray []*linkedList.LinkedList[Bucket]
	lock       sync.RWMutex
	size       int64
	seed       maphash.Seed
}

func NewSegment(initialCapacity int) *segment {
	seg := &segment{
		innerArray: make([]*linkedList.LinkedList[Bucket], initialCapacity),
		lock:       sync.RWMutex{},
		seed:       hash.MakeRandomSeed(),
	}

	return seg
}

func (h *segment) resize() {
	h.lock.Lock()
	h.seed = hash.MakeRandomSeed()
	newCapacity := 2 * len(h.innerArray)
	resizeArray := make([]*linkedList.LinkedList[Bucket], newCapacity)

	for index := range h.innerArray {
		for h.innerArray[index] != nil && h.innerArray[index].GetSize() > 0 {
			item := h.innerArray[index].RemoveFront()

			newIdx := int64(hash.Get(h.seed, item.Key) % uint64(len(resizeArray)))
			if resizeArray[newIdx] == nil {
				resizeArray[newIdx] = linkedList.NewWithInitValue[Bucket](item)
			} else {
				resizeArray[newIdx].AddFront(item)
			}
		}
	}

	h.innerArray = resizeArray
	resizeArray = nil
	h.lock.Unlock()
}

func (h *segment) checkAndStartResize() {
	h.lock.Lock()
	loadFactor := float64(h.size) / float64(len(h.innerArray))
	if loadFactor < maxLoadFactor {
		h.lock.Unlock()
		return
	}
	h.lock.Unlock()
	h.resize()
}

func (h *segment) PutSetFieldOrCreateSafe(key string, setField string) (created bool) {
	h.lock.Lock()
	hashCode := hash.Get(h.seed, key)
	index := int64(hashCode % uint64(len(h.innerArray)))
	if h.innerArray[index] != nil {
		element, found := h.innerArray[index].Find(func(current *Bucket) bool {
			return current.Key == key
		})
		if found {
			element.Value.Add(setField)
			h.lock.Unlock()
			return false
		}
	} else {
		h.innerArray[index] = linkedList.New[Bucket]()
	}

	setWithVal := set.NewSet[string]()
	setWithVal.Add(setField)
	h.innerArray[index].AddFront(&Bucket{
		Key:   key,
		Value: setWithVal,
	})
	h.size++

	h.lock.Unlock()

	h.checkAndStartResize()

	return true
}

func (h *segment) GetSafe(key string) (*set.Set[string], bool) {
	h.lock.RLock()
	defer h.lock.RUnlock()

	index := int64(hash.Get(h.seed, key) % uint64(len(h.innerArray)))
	if h.innerArray[index] == nil {
		return nil, false
	}

	element, found := h.innerArray[index].Find(func(current *Bucket) bool {
		return current.Key == key
	})
	if !found {
		return nil, false
	}

	return element.Value.Copy(), true
}

func (h *segment) RemoveSetFieldSafe(key string, setField string) (bucketRemoved bool) {
	h.lock.Lock()
	defer h.lock.Unlock()

	hashCode := hash.Get(h.seed, key)
	index := int64(hashCode % uint64(len(h.innerArray)))
	if h.innerArray[index] == nil {
		return false
	}

	element, found := h.innerArray[index].Find(func(current *Bucket) bool {
		return current.Key == key
	})
	if found {
		element.Value.Remove(setField)
		if element.Value.IsEmpty() {
			h.innerArray[index].Remove(func(current *Bucket) bool {
				return current.Key == key
			})
			h.size--
			return true
		}
	}

	return false
}

func (h *segment) GetSize() int64 {
	h.lock.RLock()
	defer h.lock.RUnlock()
	return h.size
}
