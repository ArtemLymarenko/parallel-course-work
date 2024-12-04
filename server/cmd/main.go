package main

import (
	"log"
	"parallel-course-work/pkg/threadpool"
	tcpServer "parallel-course-work/server/internal/infrastructure/tcp_server"
	tcpRouter "parallel-course-work/server/internal/infrastructure/tcp_server/router"
)

func main() {
	threadPool := threadpool.New(4, 1)
	router := tcpRouter.New()
	router.AddRoute(tcpRouter.POST, "status", func(ctx *tcpRouter.RequestContext) error {
		return ctx.ResponseJSON(tcpRouter.OK, "hello world!")
	})
	server := tcpServer.New(8080, threadPool, router)
	err := server.Start()
	if err != nil {
		log.Fatal(err)
	}
}
