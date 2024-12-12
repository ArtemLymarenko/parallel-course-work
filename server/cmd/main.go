package main

import (
	"log"
	"parallel-course-work/pkg/threadpool"
	fileReader "parallel-course-work/server/internal/infrastructure/file_reader"
	invertedIdx "parallel-course-work/server/internal/infrastructure/inverted_idx"
	tcpServer "parallel-course-work/server/internal/infrastructure/tcp_server"
	"parallel-course-work/server/internal/inteface/rest/handlers"
	v1Router "parallel-course-work/server/internal/inteface/rest/router"
)

func main() {
	threadPool := threadpool.New(4, 1)

	const resourceDir = "resources/"
	reader := fileReader.New()
	invIndex := invertedIdx.New(resourceDir, reader)

	invIndexHandlers := handlers.NewInvertedIndex(invIndex)
	router := v1Router.MustInitRouter(invIndexHandlers)

	server := tcpServer.New(8080, threadPool, router)
	err := server.Start()
	if err != nil {
		log.Fatal(err)
	}
}
