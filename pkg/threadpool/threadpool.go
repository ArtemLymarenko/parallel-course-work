package threadpool

import (
	"errors"
	"fmt"
	"log"
	"parallel-course-work/pkg/priorityqueue"
	"sync"
)

type TaskPriorityQueue interface {
	Size() int
	Push(element *Task)
	Pop() *Task
	GetItems() []*Task
	Empty() bool
}

type SyncPrimitives struct {
	mainWaiter *sync.Cond
	commonLock sync.RWMutex
	wg         sync.WaitGroup
}

type Logger interface {
	Log(...interface{})
}

type ThreadPool struct {
	logger        Logger
	mainTaskQueue TaskPriorityQueue
	sync          *SyncPrimitives
	isInitialized bool
	isTerminated  bool
}

func New(logger Logger) *ThreadPool {
	compareFunc := func(a, b *Task) bool {
		return a.CreatedAt.After(b.CreatedAt)
	}

	sp := &SyncPrimitives{
		commonLock: sync.RWMutex{},
		wg:         sync.WaitGroup{},
	}

	sp.mainWaiter = sync.NewCond(&sp.commonLock)

	return &ThreadPool{
		logger:        logger,
		mainTaskQueue: priorityqueue.New(compareFunc),
		sync:          sp,
		isInitialized: false,
		isTerminated:  false,
	}
}

func (threadPool *ThreadPool) IsWorkingUnsafe() bool {
	return threadPool.isInitialized && !threadPool.isTerminated
}

func (threadPool *ThreadPool) IsWorking() bool {
	threadPool.sync.commonLock.RLock()
	defer threadPool.sync.commonLock.RUnlock()
	return threadPool.IsWorkingUnsafe()
}

func (threadPool *ThreadPool) MustRun(mainThreadCount int) {
	threadPool.sync.commonLock.Lock()
	defer threadPool.sync.commonLock.Unlock()

	if threadPool.isInitialized || threadPool.isTerminated {
		log.Fatal("thread pool is already initialized or terminated")
	}

	for range mainThreadCount {
		threadPool.sync.wg.Add(1)
		go threadPool.routineThread()
	}

	threadPool.isInitialized = true
	threadPool.logger.Log("thread pool is running...")
}

func (threadPool *ThreadPool) MustTerminate() {
	defer threadPool.sync.wg.Wait()

	threadPool.sync.commonLock.Lock()
	defer threadPool.sync.commonLock.Unlock()

	if !threadPool.isInitialized || threadPool.isTerminated {
		return
	}

	threadPool.sync.mainWaiter.Broadcast()

	threadPool.isInitialized = false
	threadPool.isTerminated = true

	threadPool.logger.Log("threadPool terminated...")
}

func (threadPool *ThreadPool) AddTask(task *Task) error {
	if !threadPool.IsWorking() {
		return errors.New("task was not added")
	}

	threadPool.sync.commonLock.Lock()
	defer threadPool.sync.commonLock.Unlock()

	threadPool.mainTaskQueue.Push(task)
	threadPool.sync.mainWaiter.Signal()

	return nil
}

func (threadPool *ThreadPool) routineThread() {
	defer threadPool.sync.wg.Done()

	for threadPool.IsWorking() {
		task := threadPool.getTaskFromQueue()
		if task == nil {
			continue
		}

		timeTaken, err := task.Run()

		if err != nil {
			msg := fmt.Sprintf("task [%v] failed with error: %v", task.Id, err.Error())
			threadPool.logger.Log(msg)
		}

		msg := fmt.Sprintf("task [%v], finished in %v", task.Id, timeTaken)
		threadPool.logger.Log(msg)
	}
}

func (threadPool *ThreadPool) getTaskFromQueue() *Task {
	queue, waiter := threadPool.mainTaskQueue, threadPool.sync.mainWaiter

	threadPool.sync.commonLock.Lock()
	defer threadPool.sync.commonLock.Unlock()

	for queue.Empty() && !threadPool.isTerminated {
		waiter.Wait()
	}

	for {
		task := queue.Pop()
		if task == nil {
			return nil
		}

		if task != nil && task.Status == IDLE {
			_ = task.SetStatus(PROCESSING)
			msg := fmt.Sprintf("task [%v] was taken", task.Id)
			threadPool.logger.Log(msg)
			return task
		}
	}
}
