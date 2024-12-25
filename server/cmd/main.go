package main

import (
	"parallel-course-work/pkg/threadpool"
	"parallel-course-work/server/internal/app"
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
	logs := logger.MustGet("resources/logs/logs.txt", app.EnvDev)
	defer logs.Close()

	//logs := mock.NewLogger()
	threadPool := threadpool.New(logs)

	fManager := fileManager.New(logs)
	invIndex := invertedIdx.New(fManager, logs)

	const resourceDir = "resources/test/"
	invIndex.Build(resourceDir, 4)
	invIdxSchedulerService := service.NewSchedulerService(invIndex, fManager, logs)
	go invIdxSchedulerService.MonitorDirAsync(resourceDir, 30*time.Second)

	invIndexHandlers := handlers.NewInvertedIndex(invIndex, logs)
	router := v1Router.MustInitRouter(invIndexHandlers, logs)

	server := tcpServer.New(8080, threadPool, router, logs)

	threads := 12
	if err := server.Start(threads); err != nil {
		logs.Log("Server stopped with error:", err)
	}
}
