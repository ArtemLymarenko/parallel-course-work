package invertedIdx

import (
	"parallel-course-work/pkg/hash"
	"parallel-course-work/pkg/set"
	"sync"
)

const MaxSegments = 32

type SyncHashMap struct {
	syncMap       []map[string]*set.Set[string]
	locks         []sync.RWMutex
	commonLock    sync.RWMutex
	size          int
	segmentsCount int
}

func NewSyncHashMap(segmentsCount int) *SyncHashMap {
	if segmentsCount > MaxSegments {
		segmentsCount = MaxSegments
	}

	locks := make([]sync.RWMutex, segmentsCount)
	sMap := make([]map[string]*set.Set[string], segmentsCount)
	for i := 0; i < segmentsCount; i++ {
		locks[i] = sync.RWMutex{}
		sMap[i] = make(map[string]*set.Set[string])
	}

	return &SyncHashMap{
		syncMap:       sMap,
		locks:         locks,
		commonLock:    sync.RWMutex{},
		size:          0,
		segmentsCount: segmentsCount,
	}
}

func (h *SyncHashMap) Put(key string, field string) {
	hashCode := hash.GetDefault(key)
	idx := int(hashCode % uint64(h.segmentsCount))
	h.locks[idx].Lock()
	defer h.locks[idx].Unlock()
	if fileSet, ok := h.syncMap[idx][key]; ok {
		fileSet.Add(field)
		return
	}
	newSet := set.NewSet[string]()
	newSet.Add(field)
	h.syncMap[idx][key] = newSet
}

func (h *SyncHashMap) Get(key string) (*set.Set[string], bool) {
	hashCode := hash.GetDefault(key)
	idx := int(hashCode % uint64(h.segmentsCount))
	h.locks[idx].RLock()
	defer h.locks[idx].RUnlock()
	if fileSet, ok := h.syncMap[idx][key]; ok {
		return fileSet.Copy(), true
	}
	return nil, false
}

func (h *SyncHashMap) Remove(key string, field string) {
	hashCode := hash.GetDefault(key)
	idx := int(hashCode % uint64(h.segmentsCount))
	h.locks[idx].Lock()
	defer h.locks[idx].Unlock()
	if fileSet, ok := h.syncMap[idx][key]; ok {
		fileSet.Remove(field)
		if fileSet.IsEmpty() {
			delete(h.syncMap[idx], key)
		}
	}
}

func (h *SyncHashMap) GetSize() int {
	h.commonLock.RLock()
	defer h.commonLock.RUnlock()
	size := 0
	for i, _ := range h.syncMap {
		size += len(h.syncMap[i])
	}
	return size
}
