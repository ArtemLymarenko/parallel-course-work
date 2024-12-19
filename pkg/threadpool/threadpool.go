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
	mainWaiter      *sync.Cond
	secondaryWaiter *sync.Cond
	commonLock      sync.RWMutex
	printLock       sync.RWMutex
	wg              sync.WaitGroup
}

type Logger interface {
	Log(...interface{})
}

type ThreadPool struct {
	logger Logger

	mainTaskQueue        TaskPriorityQueue
	secondaryTaskQueue   TaskPriorityQueue
	sync                 *SyncPrimitives
	mainThreadCount      int
	secondaryThreadCount int
	isInitialized        bool
	isTerminated         bool
}

func New(mainThreadCount, secondaryThreadCount int, logger Logger) *ThreadPool {
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
		logger:               logger,
		mainThreadCount:      mainThreadCount,
		secondaryThreadCount: secondaryThreadCount,
		mainTaskQueue:        priorityqueue.New(compareFunc),
		secondaryTaskQueue:   priorityqueue.New(compareFunc),
		sync:                 sp,
		isInitialized:        false,
		isTerminated:         false,
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

func (threadPool *ThreadPool) MustRun() {
	threadPool.sync.commonLock.Lock()
	defer threadPool.sync.commonLock.Unlock()

	if threadPool.isInitialized || threadPool.isTerminated {
		log.Fatal("thread pool is already initialized or terminated")
	}

	for range threadPool.mainThreadCount {
		threadPool.sync.wg.Add(1)
		go threadPool.routineThread(true)
	}

	for range threadPool.secondaryThreadCount {
		threadPool.sync.wg.Add(1)
		go threadPool.routineThread(false)
	}

	threadPool.isInitialized = true
	//log.Printf("thread pool is running...\n")
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
	threadPool.sync.secondaryWaiter.Broadcast()

	threadPool.isInitialized = false
	threadPool.isTerminated = true

	//log.Printf("threadPool terminated\n")
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

func (threadPool *ThreadPool) routineThread(isPrimary bool) {
	defer threadPool.sync.wg.Done()

	for threadPool.IsWorking() {
		//Add cron job to observe main queue and move old tasks to secondary in separate thread.
		threadPool.removeOldTasks()

		task := threadPool.getTaskFromQueue(isPrimary)
		if task == nil {
			continue
		}

		timeTaken, err := task.Run()

		{
			threadPool.sync.commonLock.Lock()
			_ = task.SetStatus(FINISHED)
			threadPool.sync.commonLock.Unlock()
		}

		//threadPool.sync.printLock.Lock()
		if err != nil {
			//log.Printf("task [%v] failed with error: %v\n", task.Id, err.Error())
			msg := fmt.Sprintf("task [%v] failed with error: %v", task.Id, err.Error())
			threadPool.logger.Log(msg)
		}

		//log.Printf("task [%v], finished in %v, by primary threads: %v\n", task.Id, timeTaken, isPrimary)
		//threadPool.sync.printLock.Unlock()

		msg := fmt.Sprintf("task [%v], finished in %v, by primary threads: %v", task.Id, timeTaken, isPrimary)
		threadPool.logger.Log(msg)
	}
}

func (threadPool *ThreadPool) getQueueWithWaiter(isPrimary bool) (TaskPriorityQueue, *sync.Cond) {
	threadPool.sync.commonLock.RLock()
	defer threadPool.sync.commonLock.RUnlock()

	if isPrimary {
		return threadPool.mainTaskQueue, threadPool.sync.mainWaiter
	}

	return threadPool.secondaryTaskQueue, threadPool.sync.secondaryWaiter
}

func (threadPool *ThreadPool) getTaskFromQueue(isPrimary bool) *Task {
	queue, waiter := threadPool.getQueueWithWaiter(isPrimary)

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
			//log.Printf("task [%v] was taken\n", task.Id)
			msg := fmt.Sprintf("task [%v] was taken", task.Id)
			threadPool.logger.Log(msg)
			return task
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

			//log.Printf("task [%v] was moved\n", task.Id)
			msg := fmt.Sprintf("task [%v] was moved", task.Id)
			threadPool.logger.Log(msg)
			threadPool.sync.secondaryWaiter.Signal()
		}
	}
}
