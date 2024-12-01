package main

import (
	"parallel-course-work/pkg/threadpool"
	"parallel-course-work/server/internal/app"
)

func main() {
	tp := threadpool.New(4, 1)
	server := app.New(8080, tp)
	server.Start()
}
