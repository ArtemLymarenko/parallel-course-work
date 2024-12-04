package tcpRouter

import "errors"

var ErrInvalidRequestStatus = errors.New("invalid request status")

type ResponseStatus int

const (
	OK ResponseStatus = iota
	NotFound
	Processing
	InternalServerError
)

func (requestStatus ResponseStatus) Validate() error {
	switch requestStatus {
	case OK, NotFound, Processing:
		return nil
	default:
		return ErrInvalidRequestStatus
	}
}
