package hash

import (
	"hash/fnv"
)

func Calculate(key string) (uint64, error) {
	totalHash := fnv.New64a()
	_, err := totalHash.Write([]byte(key))
	return totalHash.Sum64(), err
}
