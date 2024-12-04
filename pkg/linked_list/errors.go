package linkedList

import "errors"

var (
	ErrElementNotFound     = errors.New("element not found")
	ErrIndexOutOfRange     = errors.New("index is out of range")
	ErrorElementNotRemoved = errors.New("element was not removed")
)
