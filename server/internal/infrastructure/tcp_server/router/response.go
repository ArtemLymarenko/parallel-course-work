package tcpRouter

import (
	"errors"
)

var ErrInvalidRequestStatus = errors.New("invalid request status")

type Response struct {
	Status ResponseStatus `json:"status"`
	Body   any            `json:"body"`
}

type ResponseStatus int

const (
	StatusOK ResponseStatus = iota
	StatusProcessing
	StatusNotFound
	StatusBadRequest
	StatusInternalServerError
)

func (responseStatus ResponseStatus) Validate() error {
	switch responseStatus {
	case StatusOK, StatusNotFound, StatusProcessing, StatusInternalServerError, StatusBadRequest:
		return nil
	default:
		return ErrInvalidRequestStatus
	}
}

func (responseStatus ResponseStatus) String() string {
	return [...]string{
		"OK",
		"Processing",
		"Not Found",
		"Bad Request",
		"Internal Server Error",
	}[responseStatus]
}
