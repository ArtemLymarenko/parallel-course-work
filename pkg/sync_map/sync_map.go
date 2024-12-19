package syncMap

import (
	"parallel-course-work/pkg/hash"
	"sync/atomic"
)

const MaxSegments = 32

type syncMap[V any] struct {
	segments      []*segment[V]
	size          atomic.Int64
	totalSegments int
}

func NewSyncHashMap[V any](initialCapacity int, segmentsCount int) *syncMap[V] {
	if segmentsCount > MaxSegments {
		segmentsCount = MaxSegments
	}

	segments := make([]*segment[V], segmentsCount)
	for i := 0; i < segmentsCount; i++ {
		segments[i] = NewSegment[V](initialCapacity)
	}

	return &syncMap[V]{
		segments:      segments,
		totalSegments: segmentsCount,
	}
}

func (h *syncMap[V]) getHashWithIndex(key string) int {
	hashCode := hash.GetDefault(key)
	return int(hashCode % uint64(h.totalSegments))
}

func (h *syncMap[V]) Put(key string, value V) {
	bucket := &Bucket[V]{
		Key:   key,
		Value: value,
	}

	idx := h.getHashWithIndex(key)
	exists := h.segments[idx].PutSafe(bucket)
	if !exists {
		h.size.Add(1)
	}
}

func (h *syncMap[V]) Get(key string) (res V, ok bool) {
	idx := h.getHashWithIndex(key)
	bucket, ok := h.segments[idx].GetSafe(key)
	if !ok {
		return res, false
	}
	return bucket.Value, true
}

func (h *syncMap[V]) Remove(key string) {
	idx := h.getHashWithIndex(key)
	h.segments[idx].RemoveSafe(key)
	h.size.Add(-1)
}

func (h *syncMap[V]) Modify(key string, cb func(modify V) interface{}) (bool, interface{}) {
	idx := h.getHashWithIndex(key)
	return h.segments[idx].ModifySafe(key, cb)
}

func (h *syncMap[V]) GetSize() int64 {
	return h.size.Load()
}
