package hash

import (
	"hash/fnv"
	"hash/maphash"
)

func MakeRandomSeed() maphash.Seed {
	return maphash.MakeSeed()
}

func Get(seed maphash.Seed, key string) uint64 {
	var h maphash.Hash
	h.SetSeed(seed)
	_, _ = h.WriteString(key)
	return h.Sum64()
}

func GetDefault(key string) uint64 {
	h := fnv.New64()
	_, _ = h.Write([]byte(key))
	return h.Sum64()
}
