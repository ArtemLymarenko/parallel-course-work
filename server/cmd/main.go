package main

import (
	"parallel-course-work/pkg/threadpool"
	"sync/atomic"
	"time"
)

func main() {
	pool := threadpool.New(10, 1)
	pool.MustRun()

	id := atomic.Int64{}
	id.Store(0)

	for range 100 {
		task := threadpool.NewTask(id.Add(1), func() error {
			time.Sleep(1 * time.Second)
			return nil
		})

		pool.AddTask(task)
	}
	time.Sleep(6 * time.Second)
	pool.MustTerminate()
}
