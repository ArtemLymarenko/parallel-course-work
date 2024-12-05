package hash

import (
	"errors"
	"hash/fnv"
	"strings"
)

var ErrCalculatingHash = errors.New("error calculating hash")

func Calculate(key string) (uint64, error) {
	if len(strings.Trim(key, " ")) == 0 {
		return 0, ErrCalculatingHash
	}

	totalHash := fnv.New64a()
	_, err := totalHash.Write([]byte(key))
	if err != nil {
		return 0, ErrCalculatingHash
	}

	return totalHash.Sum64(), nil
}
