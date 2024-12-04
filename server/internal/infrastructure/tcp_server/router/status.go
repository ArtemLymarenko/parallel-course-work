package tcpRouter

import "errors"

var ErrInvalidRequestStatus = errors.New("invalid request status")

type ResponseStatus int

const (
	OK ResponseStatus = iota
	Processing
	DataNotFound
	InternalServerError
)

func (requestStatus ResponseStatus) Validate() error {
	switch requestStatus {
	case OK, DataNotFound, Processing, InternalServerError:
		return nil
	default:
		return ErrInvalidRequestStatus
	}
}
