package tcpClient

import (
	"encoding/json"
	"errors"
)

var ErrInvalidRequestStatus = errors.New("invalid request status")

type ResponseStatus int

const (
	StatusOK ResponseStatus = iota
	StatusProcessing
	StatusNotFound
	StatusBadRequest
	StatusInternalServerError
)

func (requestStatus ResponseStatus) Validate() error {
	switch requestStatus {
	case StatusOK, StatusNotFound, StatusProcessing, StatusInternalServerError, StatusBadRequest:
		return nil
	default:
		return ErrInvalidRequestStatus
	}
}

func (requestStatus ResponseStatus) String() string {
	return [...]string{
		"OK",
		"Processing",
		"Not Found",
		"Bad Request",
		"Internal Server Error",
	}[requestStatus]
}

type Response struct {
	Status ResponseStatus  `json:"status"`
	Body   json.RawMessage `json:"body"`
}
