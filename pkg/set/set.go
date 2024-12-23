package set

type Set[K comparable] struct {
	Keys map[K]struct{}
}

func NewSet[K comparable]() *Set[K] {
	return &Set[K]{
		Keys: make(map[K]struct{}),
	}
}

func (set *Set[K]) Add(key K) {
	set.Keys[key] = struct{}{}
}

func (set *Set[K]) Has(key K) bool {
	_, ok := set.Keys[key]
	return ok
}

func (set *Set[K]) Remove(key K) {
	delete(set.Keys, key)
}

func (set *Set[K]) IsEmpty() bool {
	return len(set.Keys) == 0
}

func (set *Set[K]) ToSlice() []K {
	keys := make([]K, len(set.Keys))
	i := 0
	for key := range set.Keys {
		keys[i] = key
		i++
	}
	return keys
}

func (set *Set[K]) Copy() *Set[K] {
	copySet := NewSet[K]()
	for key := range set.Keys {
		copySet.Add(key)
	}
	return copySet
}
