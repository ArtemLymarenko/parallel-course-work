package threadpool

import "errors"

var (
	ErrInvalidTaskStatus = errors.New("invalid task status")
)

type TaskStatus int

const (
	IDLE TaskStatus = iota
	PROCESSING
	FINISHED
)

func (taskStatus TaskStatus) Validate() error {
	switch taskStatus {
	case IDLE, PROCESSING, FINISHED:
		return nil
	default:
		return ErrInvalidTaskStatus
	}
}
