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
	IsMoved   bool
}

func NewTask(id int64, runFunc RunFunc) *Task {
	return &Task{
		Id:        id,
		CreatedAt: time.Now(),
		Status:    IDLE,
		RunFunc:   runFunc,
		IsMoved:   false,
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

func (task *Task) SetMoved(moved bool) {
	task.IsMoved = moved
}

func (task *Task) IsOld() bool {
	return task.Status == IDLE && !task.IsMoved && time.Since(task.CreatedAt) > 5*time.Second
}
