package syncMap

//
//import (
//	"parallel-course-work/pkg/hash"
//)
//
//const (
//	SegmentsLen = 16
//)
//
//type syncHashMap[V any] struct {
//	segments []*Segment[V]
//}
//
//func NewSyncHashMap[V any](initialCapacity int, loadFactor float64) *syncHashMap[V] {
//	segments := make([]*Segment[V], SegmentsLen)
//	for i := 0; i < SegmentsLen; i++ {
//		segments[i] = NewSegment[V](initialCapacity, loadFactor)
//	}
//
//	return &syncHashMap[V]{
//		segments: segments,
//	}
//}
//
//func (h *syncHashMap[V]) getSegmentIndexFromKey(key string) (int, error) {
//	hashCode, err := hash.Calculate(key)
//	if err != nil {
//		return 0, err
//	}
//
//	return int(hashCode) % SegmentsLen, nil
//}
//
//func (h *syncHashMap[V]) Put(key string, value V) error {
//	idx, err := h.getSegmentIndexFromKey(key)
//	if err != nil {
//		return err
//	}
//
//	err = h.segments[idx].PutSafe(key, value)
//	return err
//}
//
//func (h *syncHashMap[V]) Get(key string) (*Bucket[V], bool) {
//	idx, err := h.getSegmentIndexFromKey(key)
//	if err != nil {
//		return nil, false
//	}
//
//	return h.segments[idx].GetSafe(key)
//}
//
//func (h *syncHashMap[V]) Remove(key string) error {
//	idx, err := h.getSegmentIndexFromKey(key)
//	if err != nil {
//		return err
//	}
//
//	return h.segments[idx].RemoveSafe(key)
//}
