package linkedList

import (
	"math/rand"
	"testing"
)

type ScholarShip struct {
	cardNumber  int64
	amount      int
	isAvailable bool
}

func Test(t *testing.T) {
	list := New[ScholarShip]()
	const (
		maxPushFront   = 50_000
		maxPushBack    = 10_000
		maxRead        = 20_000
		maxRemoveFront = 5_000
		maxRemoveBack  = 5_000
	)
	var i = 0
	for i < maxPushFront {
		amount := rand.Intn(2000) + 1000
		cardNumber := rand.Int63()
		list.AddFront(ScholarShip{cardNumber: cardNumber, amount: amount, isAvailable: true})
		i++
	}
	if list.GetSize() != maxPushFront {
		t.Errorf("Expected %d, got %d", maxPushFront, list.GetSize())
	}

	i = 0
	for i < maxPushBack {
		amount := rand.Intn(2000) + 1000
		cardNumber := rand.Int63()
		list.AddBack(ScholarShip{cardNumber: cardNumber, amount: amount, isAvailable: true})
		i++
	}
	if list.GetSize() != maxPushFront+maxPushBack {
		t.Errorf("Expected %d, got %d", maxPushFront+maxPushBack, list.GetSize())
	}

	i = 0
	for i < maxRead {
		index := rand.Intn(list.GetSize())
		_, err := list.FindByIndex(index)
		if err != nil {
			t.Log(err)
		}
		i++
	}
	if list.GetSize() != maxPushFront+maxPushBack {
		t.Errorf("Expected %d, got %d", maxPushFront+maxPushBack, list.GetSize())
	}

	i = 0
	for i < maxRemoveFront {
		list.RemoveByIndex(0)
		i++
	}
	if list.GetSize() != maxPushFront+maxPushBack-maxRemoveFront {
		t.Errorf("Expected %d, got %d", maxPushFront+maxPushBack-maxRemoveFront, list.GetSize())
	}

	i = 0
	for i < maxRemoveBack {
		list.RemoveByIndex(list.GetSize() - 1)
		i++
	}
}
