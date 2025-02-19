package threadpool

import (
	"time"
)

type RunFunc func() error

type Task struct {
	Id        int64
	CreatedAt time.Time
	Status    TaskStatus
	RunFunc   RunFunc
}

func NewTask(id int64, runFunc RunFunc) *Task {
	return &Task{
		Id:        id,
		CreatedAt: time.Now(),
		Status:    IDLE,
		RunFunc:   runFunc,
	}
}

func (task *Task) Run() (time.Duration, error) {
	now := time.Now()
	err := task.RunFunc()
	return time.Since(now), err
}

func (task *Task) SetRunFunc(runFunc RunFunc) {
	task.RunFunc = runFunc
}

func (task *Task) SetId(id int64) {
	task.Id = id
}

func (task *Task) SetCreatedAt(createdAt time.Time) {
	task.CreatedAt = createdAt
}

func (task *Task) SetStatus(status TaskStatus) error {
	if err := status.Validate(); err != nil {
		return err
	}

	task.Status = status
	return nil
}
