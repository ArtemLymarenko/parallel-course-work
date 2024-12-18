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

func (h *syncMap[V]) getHashWithIndex(key string) (uint64, int) {
	hashCode, _ := hash.Calculate(key)
	return hashCode, int(hashCode % uint64(h.totalSegments))
}

func (h *syncMap[V]) Put(key string, value V) {
	hashCode, idx := h.getHashWithIndex(key)
	bucket := &Bucket[V]{
		Key:   key,
		Value: value,
		hash:  hashCode,
	}

	h.segments[idx].PutSafe(bucket)
	h.size.Add(1)
}

func (h *syncMap[V]) Get(key string) (res V, ok bool) {
	hashCode, idx := h.getHashWithIndex(key)
	bucket, ok := h.segments[idx].GetSafe(hashCode, key)
	if !ok {
		return res, false
	}
	return bucket.Value, true
}

func (h *syncMap[V]) Remove(key string) {
	hashCode, idx := h.getHashWithIndex(key)
	h.segments[idx].RemoveSafe(hashCode, key)
	h.size.Add(-1)
}

func (h *syncMap[V]) GetSize() int64 {
	return h.size.Load()
}
