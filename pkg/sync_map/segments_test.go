package syncMap

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestSegment_PutSafe(t *testing.T) {
	segment := New[int](3)
	go func() {
		segment.PutSafe("keasdasdasdy", 1) //Same hash
	}()
	go func() {
		segment.PutSafe("kasdfasdfeey", 1) //Same hash
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
	repeatCount := 100

	for i := 0; i < repeatCount; i++ {
		t.Run(fmt.Sprintf("TestRun %d", i+1), func(t *testing.T) {
			initCap := 16
			expectedSize := int64(5000)
			addElements := 5000

			segment := New[int](initCap)
			wg := sync.WaitGroup{}
			wg.Add(addElements)

			for i := 0; i < addElements; i++ {
				go func(i int) {
					segment.PutSafe(strconv.Itoa(i), 1)
					wg.Done()
				}(i)
			}

			wg.Wait()

			// Перевірка розміру після ресайзу
			if segment.GetSize() != expectedSize {
				t.Errorf("expected size to be %d, got %d", expectedSize, len(segment.innerArray))
			}

			// Перевірка емності сегменту після ресайзу
			if len(segment.innerArray) < addElements {
				t.Errorf("expected cap to be %d, got %d", len(segment.innerArray), initCap)
			}

			// Перевірка, що resizeArray порожній
			if segment.resizeArray != nil {
				t.Errorf("expected resizeArray to be nil, got %v", segment.resizeArray)
			}

			notFound := 0
			for i := 0; i < addElements; i++ {
				_, ok := segment.GetSafe(strconv.Itoa(i))
				if !ok {
					notFound++
					t.Errorf("expected to find this element %v", strconv.Itoa(i))
				}
			}

			if notFound != 0 {
				t.Errorf("expected to find all, not found %v", notFound)
			}
		})
	}
}
