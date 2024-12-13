package tcpRouter

import "errors"

var ErrInvalidRequestStatus = errors.New("invalid request status")

type ResponseStatus int

const (
	StatusOK ResponseStatus = iota + 1
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
