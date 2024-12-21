package invertedIdx

import (
	"parallel-course-work/pkg/hash"
	"parallel-course-work/pkg/set"
	"sync/atomic"
)

const MaxSegments = 32

type syncMap struct {
	segments      []*segment
	size          atomic.Int64
	totalSegments int
}

func NewSyncHashMap(initialCapacity int, segmentsCount int) *syncMap {
	if segmentsCount > MaxSegments {
		segmentsCount = MaxSegments
	}

	segments := make([]*segment, segmentsCount)
	for i := 0; i < segmentsCount; i++ {
		segments[i] = NewSegment(initialCapacity)
	}

	return &syncMap{
		segments:      segments,
		totalSegments: segmentsCount,
	}
}

func (h *syncMap) getHashWithIndex(key string) int {
	hashCode := hash.GetDefault(key)
	return int(hashCode % uint64(h.totalSegments))
}

func (h *syncMap) Put(key string, field string) {
	idx := h.getHashWithIndex(key)
	created := h.segments[idx].PutSetFieldOrCreateSafe(key, field)
	if created {
		h.size.Add(1)
	}
}

func (h *syncMap) Get(key string) (*set.Set[string], bool) {
	idx := h.getHashWithIndex(key)
	result, ok := h.segments[idx].GetSafe(key)
	if !ok {
		return nil, false
	}
	return result, true
}

func (h *syncMap) Remove(key string, field string) {
	idx := h.getHashWithIndex(key)
	bucketRemoved := h.segments[idx].RemoveSetFieldSafe(key, field)
	if bucketRemoved {
		h.size.Add(-1)
	}
}

func (h *syncMap) GetSize() int64 {
	return h.size.Load()
}
