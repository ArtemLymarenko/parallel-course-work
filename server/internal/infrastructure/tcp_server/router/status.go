package tcpRouter

import "errors"

var ErrInvalidRequestStatus = errors.New("invalid request status")

type ResponseStatus int

const (
	OK ResponseStatus = iota
	Processing
	NotFound
	BadRequest
	InternalServerError
)

func (requestStatus ResponseStatus) Validate() error {
	switch requestStatus {
	case OK, NotFound, Processing, InternalServerError, BadRequest:
		return nil
	default:
		return ErrInvalidRequestStatus
	}
}
