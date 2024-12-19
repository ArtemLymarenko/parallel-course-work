package threadpool

import (
	"fmt"
	"parallel-course-work/pkg/mock"
	"testing"
	"time"
)

var logs = mock.NewLogger()

func TestThreadPoolTasks(t *testing.T) {
	pool := New(1, 1, logs)
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

func TestRunPrimaryAndSecondaryAtTime(t *testing.T) {
	const numIterations = 5

	for i := 0; i < numIterations; i++ {
		t.Run(fmt.Sprintf("Iteration-%d", i+1), func(t *testing.T) {
			pool := New(1, 1, logs)

			taskExecuted := false
			task := NewTask(1, func() error {
				if taskExecuted {
					t.Errorf("task was executed more than once")
				}
				taskExecuted = true
				time.Sleep(1 * time.Second)
				return nil
			})

			pool.mainTaskQueue.Push(task)
			pool.secondaryTaskQueue.Push(task)

			pool.MustRun()
			time.Sleep(2 * time.Second)
			if pool.mainTaskQueue.Size() != 0 {
				t.Errorf("expected main queue to be empty, but got size %d", pool.mainTaskQueue.Size())
			}

			if pool.secondaryTaskQueue.Size() != 0 {
				t.Errorf("expected secondary queue to be empty, but got size %d", pool.secondaryTaskQueue.Size())
			}

			pool.MustTerminate()
		})
	}
}
