package syncMap

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestSyncMap_Put(t *testing.T) {
	m := NewSyncHashMapV2[string](3, 2)
	go func() {
		m.Put("keasdasdasdy", "") //Same hash
	}()
	go func() {
		m.Put("kasdfasdfeey", "") //Same hash
	}()
	time.Sleep(1 * time.Second)

	_, ok := m.Get("keasdasdasdy")
	if !ok {
		t.Error("expected key to be present")
	}

	_, ok = m.Get("kasdfasdfeey")
	if !ok {
		t.Error("expected key to be present")
	}
}

func TestSyncMap_Resize(t *testing.T) {
	const iterations = 100
	for iter := 0; iter < iterations; iter++ {
		t.Run(fmt.Sprintf("Iteration #%d", iter+1), func(t *testing.T) {
			initCap := 120
			expectedSize := int64(50000)
			addElements := 50000

			m := NewSyncHashMapV2[string](initCap, 32)
			wg := sync.WaitGroup{}
			wg.Add(addElements)

			for i := 0; i < addElements; i++ {
				go func(i int) {
					m.Put(strconv.Itoa(i), "")
					wg.Done()
				}(i)
			}

			wg.Wait()

			if m.GetSize() != expectedSize {
				t.Errorf("expected size to be %d, got %d", expectedSize, m.GetSize())
			}

			notFound := 0
			for i := 0; i < addElements; i++ {
				_, ok := m.Get(strconv.Itoa(i))
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

func TestSyncMap_Remove(t *testing.T) {
	initCap := 16
	expectedSize := int64(50000)
	addElements := 50000

	m := NewSyncHashMapV2[string](initCap, 32)
	wg := sync.WaitGroup{}
	wg.Add(addElements)

	for i := 0; i < addElements; i++ {
		go func(i int) {
			m.Put(strconv.Itoa(i), "")
			wg.Done()
		}(i)
	}

	wg.Wait()

	for i := 0; i < addElements; i++ {
		m.Remove(strconv.Itoa(i))
	}

	if m.GetSize() != 0 {
		t.Errorf("expected size to be %d, got %d", expectedSize, m.GetSize())
	}
}
