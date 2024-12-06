package syncMap

import (
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestSegment_PutSafe(t *testing.T) {
	segment := NewSegment[int](3, 0.75)
	go func() {
		_ = segment.PutSafe("keasdasdasdy", 1) //Same hash
	}()
	go func() {
		_ = segment.PutSafe("kasdfasdfeey", 1) //Same hash
	}()
	time.Sleep(1 * time.Second)
	_, ok := segment.GetSafe("keasdasdasdy")
	if !ok {
		t.Error("expected key to be present")
	}

	_, ok = segment.GetSafe("kasdfasdfeey")
	if !ok {
		t.Error("expected key to be present")
	}
}

func TestSegment_Resize(t *testing.T) {
	initCap := 10
	expectedSize := 1280
	addElements := 1000

	segment := NewSegment[int](initCap, 0.75)
	wg := sync.WaitGroup{}
	wg.Add(addElements)
	for i := range addElements {
		go func() {
			_ = segment.PutSafe(strconv.Itoa(i), 1)
			wg.Done()
		}()
	}

	wg.Wait()

	size := segment.GetInnerArrayLenSafe()
	if size != expectedSize {
		t.Errorf("expected size to be %d, got %d", expectedSize, size)
	}
}
