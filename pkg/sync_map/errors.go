package syncMap

import "errors"

var (
	ErrCalculatingHash = errors.New("error calculating hash")
	ErrElementNotFound = errors.New("element not found")
)
