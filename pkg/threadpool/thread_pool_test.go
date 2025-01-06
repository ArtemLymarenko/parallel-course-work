package threadpool

import (
	"fmt"
	"github.com/ArtemLymarenko/parallel-course-work/pkg/mock"
	"testing"
	"time"
)

var logs = mock.NewLogger()

func TestRun(t *testing.T) {
	pool := New(logs)
	for range 50000 {
		go func() {
			task := NewTask(1, func() error {
				time.Sleep(1 * time.Second)
				return nil
			})
			pool.AddTask(task)
		}()
	}

	time.Sleep(4 * time.Second)
	fmt.Println(pool.mainTaskQueue.Size())
}
