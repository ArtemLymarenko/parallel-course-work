package logger

import (
	"fmt"
	"log"
	"parallel-course-work/server/internal/app"
	"sync"
	"time"
)

type FileLogger interface {
	LogUnsafe(...interface{})
	Close()
}

type logger struct {
	fileLogger *fileLogger
	lock       sync.Mutex
	env        app.Env
}

func MustGet(path string, env app.Env) *logger {
	fl, err := NewFileLogger(path, 20)
	if err != nil || !env.Valid() {
		log.Fatal(err)
	}

	return &logger{
		fileLogger: fl,
		lock:       sync.Mutex{},
		env:        env,
	}
}

const TimeFormat = "2006/01/02 15:04:05"

func (l *logger) Log(v ...interface{}) {
	l.lock.Lock()
	defer l.lock.Unlock()
	timestamp := time.Now().Format(TimeFormat)
	msg := fmt.Sprint(v...)
	logMsg := fmt.Sprintf("%s %s", timestamp, msg)

	log.Printf(msg)
	if l.env.IsProduction() {
		l.fileLogger.LogUnsafe(logMsg)
	}
}

func (l *logger) Close() {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.fileLogger.Close()
}
