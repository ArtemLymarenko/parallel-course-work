package threadpool

import (
	"fmt"
	"log"
	"parallel-course-work/pkg/priorityqueue"
	"sync"
	"time"
)

type TaskPriorityQueue interface {
	Size() int
	Push(element *Task)
	Pop() (*Task, error)
	GetItems() []*Task
	Empty() bool
}

type State struct {
	isInitialized bool
	isTerminated  bool
	isPaused      bool
}

type SyncPrimitives struct {
	mainWaiter      *sync.Cond
	secondaryWaiter *sync.Cond
	commonLock      sync.RWMutex
	printLock       sync.RWMutex
	wg              sync.WaitGroup
}

type ThreadPool struct {
	mainThreadCount      int
	secondaryThreadCount int
	mainTaskQueue        TaskPriorityQueue
	secondaryTaskQueue   TaskPriorityQueue
	state                State
	sync                 *SyncPrimitives
}

func New(mainThreadCount, secondaryThreadCount int) *ThreadPool {
	state := State{false, false, false}
	compareFunc := func(a, b *Task) bool {
		return a.CreatedAt.After(b.CreatedAt)
	}

	sp := &SyncPrimitives{
		commonLock: sync.RWMutex{},
		printLock:  sync.RWMutex{},
		wg:         sync.WaitGroup{},
	}

	sp.mainWaiter = sync.NewCond(&sp.commonLock)
	sp.secondaryWaiter = sync.NewCond(&sp.commonLock)

	return &ThreadPool{
		mainThreadCount:      mainThreadCount,
		secondaryThreadCount: secondaryThreadCount,
		mainTaskQueue:        priorityqueue.New(compareFunc),
		secondaryTaskQueue:   priorityqueue.New(compareFunc),
		sync:                 sp,
		state:                state,
	}
}

func (threadPool *ThreadPool) IsWorkingUnsafe() bool {
	return threadPool.state.isInitialized && !threadPool.state.isTerminated && !threadPool.state.isPaused
}

func (threadPool *ThreadPool) IsWorking() bool {
	threadPool.sync.commonLock.RLock()
	defer threadPool.sync.commonLock.RUnlock()
	return threadPool.IsWorkingUnsafe()
}

func (threadPool *ThreadPool) MustRun() {
	threadPool.sync.commonLock.Lock()
	defer threadPool.sync.commonLock.Unlock()

	if threadPool.state.isInitialized || threadPool.state.isTerminated {
		log.Fatal("thread pool is already initialized or terminated")
	}

	for range threadPool.mainThreadCount {
		threadPool.sync.wg.Add(1)
		go threadPool.routine(true)
	}

	for range threadPool.secondaryThreadCount {
		threadPool.sync.wg.Add(1)
		go threadPool.routine(false)
	}

	threadPool.state.isInitialized = true
	fmt.Println("thread pool is running!")
}

func (threadPool *ThreadPool) Terminate() {
	defer threadPool.sync.wg.Wait()

	threadPool.sync.commonLock.Lock()
	defer threadPool.sync.commonLock.Unlock()

	if !threadPool.state.isInitialized || threadPool.state.isTerminated {
		return
	}

	threadPool.sync.mainWaiter.Broadcast()
	threadPool.sync.secondaryWaiter.Broadcast()

	threadPool.state.isInitialized = false
	threadPool.state.isTerminated = true

	fmt.Println("threadPool terminated")
}

func (threadPool *ThreadPool) AddTask(task *Task) bool {
	if !threadPool.IsWorking() {
		return false
	}

	threadPool.sync.commonLock.Lock()
	defer threadPool.sync.commonLock.Unlock()

	threadPool.mainTaskQueue.Push(task)
	threadPool.sync.mainWaiter.Signal()

	fmt.Printf("added new task - %v\n", task.Id)
	return true
}

func (threadPool *ThreadPool) routine(isMain bool) {
	defer threadPool.sync.wg.Done()

	for !threadPool.state.isTerminated {
		threadPool.removeOldTasks()

		task, err := threadPool.getTaskFromQueue(isMain)
		if err != nil {
			return
		}

		now := time.Now()
		if err = task.Run(); err != nil {
			fmt.Printf("task [%v] failed with error: %v\n", task.Id, err.Error())
		}

		fmt.Printf("task [%v], finished in - %v\n", task.Id, time.Since(now))
	}
}

func (threadPool *ThreadPool) getQueueWithWaiter(isMain bool) (TaskPriorityQueue, *sync.Cond) {
	threadPool.sync.commonLock.RLock()
	defer threadPool.sync.commonLock.RUnlock()

	if isMain {
		return threadPool.mainTaskQueue, threadPool.sync.mainWaiter
	}

	return threadPool.secondaryTaskQueue, threadPool.sync.secondaryWaiter
}

func (threadPool *ThreadPool) getTaskFromQueue(isMain bool) (*Task, error) {
	queue, waiter := threadPool.getQueueWithWaiter(isMain)

	{
		threadPool.sync.commonLock.Lock()
		defer threadPool.sync.commonLock.Unlock()

		for queue.Empty() && !threadPool.state.isTerminated {
			waiter.Wait()
		}

		for {
			task, err := queue.Pop()
			if err != nil {
				return nil, err
			}

			if task != nil && task.Status == IDLE {
				_ = task.SetStatus(PROCESSING)
				fmt.Printf("task [%v] was taken\n", task.Id)
				return task, nil
			}
		}
	}
}

func (threadPool *ThreadPool) removeOldTasks() {
	threadPool.sync.commonLock.Lock()
	defer threadPool.sync.commonLock.Unlock()

	for _, task := range threadPool.mainTaskQueue.GetItems() {
		if task.IsOld() {
			threadPool.secondaryTaskQueue.Push(task)
			task.SetMoved(true)

			fmt.Printf("task [%v] was moved\n", task.Id)
			threadPool.sync.secondaryWaiter.Signal()
		}
	}
}
