package threadpool

import (
	"testing"
	"time"
)

func TestThreadPoolTasks(t *testing.T) {
	pool := New(1, 1)
	task := NewTask(1, func() error {
		time.Sleep(1 * time.Second)
		return nil
	})
	pool.mainTaskQueue.Push(task)
	pool.secondaryTaskQueue.Push(task)

	pool.MustRun()

	time.Sleep(2 * time.Second)

	task2 := NewTask(2, func() error {
		time.Sleep(1 * time.Second)
		return nil
	})
	pool.secondaryTaskQueue.Push(task2)
	pool.sync.secondaryWaiter.Signal()

	time.Sleep(4 * time.Second)

	if pool.secondaryTaskQueue.Size() != 0 {
		t.Errorf("expected to finish the task %d", pool.secondaryTaskQueue.Size())
	}

	pool.MustTerminate()
}
