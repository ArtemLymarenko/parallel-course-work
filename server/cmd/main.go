package main

import (
	"github.com/ArtemLymarenko/parallel-course-work/pkg/mock"
	"github.com/ArtemLymarenko/parallel-course-work/pkg/threadpool"
	"server/internal/app"
	filemanager "server/internal/infrastructure/file_manager"
	invertedIdx "server/internal/infrastructure/inverted_idx"
	"server/internal/infrastructure/logger"
	tcpServer "server/internal/infrastructure/tcp_server"
	"server/internal/inteface/rest/handlers"
	v1Router "server/internal/inteface/rest/router"
	"server/internal/service"
	"time"
)

func main() {
	_ = mock.NewLogger()
	loggerService := logger.MustGet("resources/logs/logs.txt", app.EnvProd)
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
