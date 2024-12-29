package invertedIdx

import (
	"github.com/ArtemLymarenko/parallel-course-work/pkg/hash"
	"github.com/ArtemLymarenko/parallel-course-work/pkg/set"
	"sync"
	"sync/atomic"
)

const MaxSegments = 32

type SyncHashMap struct {
	syncMap       []map[string]*set.Set[string]
	locks         []sync.RWMutex
	segmentsCount int
	size          atomic.Int64
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
	h.size.Add(1)
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
			h.size.Add(-1)
		}
	}
}

func (h *SyncHashMap) GetSize() int64 {
	return h.size.Load()
}
