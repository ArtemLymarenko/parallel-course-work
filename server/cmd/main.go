package main

import (
	"parallel-course-work/pkg/threadpool"
	"sync/atomic"
	"time"
)

func main() {
	pool := threadpool.New(4, 1)
	pool.MustRun()

	id := atomic.Int64{}
	id.Store(0)

	for range 10 {
		task := threadpool.NewTask(id.Add(1), func() error {
			time.Sleep(9 * time.Second)
			return nil
		})

		pool.AddTask(task)
	}
	time.Sleep(10 * time.Second)
	pool.Terminate()
}
