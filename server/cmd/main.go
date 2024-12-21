package main

import (
	"parallel-course-work/pkg/threadpool"
	fileManager "parallel-course-work/server/internal/infrastructure/file_manager"
	invertedIdx "parallel-course-work/server/internal/infrastructure/inverted_idx"
	"parallel-course-work/server/internal/infrastructure/logger"
	tcpServer "parallel-course-work/server/internal/infrastructure/tcp_server"
	"parallel-course-work/server/internal/inteface/rest/handlers"
	v1Router "parallel-course-work/server/internal/inteface/rest/router"
	"parallel-course-work/server/internal/service"
	"time"
)

func main() {
	logs := logger.MustGet()
	defer logs.Close()

	threadPool := threadpool.New(4, 1, logs)

	const resourceDir = "resources/test/"
	fManager := fileManager.New(logs)
	invIndex := invertedIdx.New(resourceDir, fManager, logs)

	period := 30 * time.Second
	invIdxSchedulerService := service.NewSchedulerService(invIndex, fManager, resourceDir, period, logs)
	go invIdxSchedulerService.ScheduleAsync()

	invIndexHandlers := handlers.NewInvertedIndex(invIndex, logs)
	router := v1Router.MustInitRouter(invIndexHandlers, logs)

	server := tcpServer.New(8080, threadPool, router, logs)
	if err := server.Start(); err != nil {
		logs.Log("Server stopped with error:", err)
	}
}
