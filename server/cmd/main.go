package main

import (
	"parallel-course-work/pkg/mock"
	"parallel-course-work/pkg/threadpool"
	"parallel-course-work/server/internal/app"
	filemanager "parallel-course-work/server/internal/infrastructure/file_manager"
	invertedIdx "parallel-course-work/server/internal/infrastructure/inverted_idx"
	"parallel-course-work/server/internal/infrastructure/logger"
	tcpServer "parallel-course-work/server/internal/infrastructure/tcp_server"
	"parallel-course-work/server/internal/inteface/rest/handlers"
	v1Router "parallel-course-work/server/internal/inteface/rest/router"
	"parallel-course-work/server/internal/service"
	"time"
)

func main() {
	_ = mock.NewLogger()
	loggerService := logger.MustGet("resources/logs/logs.txt", app.EnvDev)
	defer loggerService.Close()

	fileManager := filemanager.New(loggerService)
	invIndex := invertedIdx.New(fileManager, loggerService)

	const resourceDir = "resources/data/"
	invIndex.Build(resourceDir, 12)

	invIdxSchedulerService := service.NewSchedulerService(invIndex, fileManager, loggerService)
	go invIdxSchedulerService.MonitorDirAsync(resourceDir, 30*time.Second)

	invIndexHandlers := handlers.NewInvertedIndex(invIndex, loggerService)
	router := v1Router.MustInitRouter(invIndexHandlers, loggerService)

	threadPool := threadpool.New(loggerService)
	server := tcpServer.New(8080, threadPool, router, loggerService)

	const threadCount = 12
	if err := server.Start(threadCount); err != nil {
		loggerService.Log("Server stopped with error:", err)
	}
}
