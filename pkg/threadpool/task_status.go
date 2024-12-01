package threadpool

import "errors"

var (
	ErrInvalidTaskStatus = errors.New("invalid task status")
)

type TaskStatus int

const (
	IDLE TaskStatus = iota
	PROCESSING
)

func (taskStatus TaskStatus) Validate() error {
	switch taskStatus {
	case IDLE, PROCESSING:
		return nil
	default:
		return ErrInvalidTaskStatus
	}
}
