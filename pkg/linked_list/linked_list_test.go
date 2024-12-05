package linkedList

import (
	"math/rand"
	"testing"
)

type ScholarShip struct {
	CardNumber  int64
	Amount      int
	IsAvailable bool
}

func Test(t *testing.T) {
	list := New[ScholarShip]()
	const (
		maxPushFront      = 50_000
		maxPushBack       = 10_000
		maxRemoveFront    = 20_000
		maxReadCardNumber = 5_000
		maxRemoveByName   = 5_000
	)
	for range maxPushFront {
		amount := rand.Intn(2000) + 1000
		cardNumber := rand.Int63()
		list.AddFront(&ScholarShip{CardNumber: cardNumber, Amount: amount, IsAvailable: true})
	}
	if list.GetSize() != maxPushFront {
		t.Errorf("Expected %d, got %d", maxPushFront, list.GetSize())
	}

	for range maxPushBack {
		amount := rand.Intn(2000) + 1000
		cardNumber := rand.Int63()
		list.AddBack(&ScholarShip{CardNumber: cardNumber, Amount: amount, IsAvailable: true})
	}
	if list.GetSize() != maxPushFront+maxPushBack {
		t.Errorf("Expected %d, got %d", maxPushFront+maxPushBack, list.GetSize())
	}

	for range maxRemoveFront {
		list.RemoveFront()
	}
	if list.GetSize() != maxPushFront+maxPushBack-maxRemoveFront {
		t.Errorf("Expected %d, got %d", maxPushFront+maxPushBack-maxRemoveFront, list.GetSize())
	}

	for range maxReadCardNumber {
		_, ok := list.FindByStructField("IsAvailable", true)
		if !ok {
			t.Errorf("Expected to read data, got %v", ok)
		}
	}

	for range maxRemoveByName {
		err := list.RemoveByStructField("IsAvailable", true)
		if err != nil {
			t.Errorf("Expected to read data, got %v", err)
		}
	}
	if list.GetSize() != maxPushFront+maxPushBack-maxRemoveFront-maxRemoveByName {
		t.Errorf("Expected %d, got %d", maxPushFront+maxPushBack-maxRemoveFront-maxRemoveByName, list.GetSize())
	}
}
