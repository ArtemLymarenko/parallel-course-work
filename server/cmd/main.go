package main

import (
	"parallel-course-work/pkg/threadpool"
	tcpServer "parallel-course-work/server/internal/infrastructure/tcp_server"
	tcpRouter "parallel-course-work/server/internal/infrastructure/tcp_server/router"
)

func main() {
	threadPool := threadpool.New(4, 1)
	router := tcpRouter.New()
	router.AddRoute("status", tcpRouter.POST, func(ctx *tcpRouter.RequestContext) error {
		return ctx.ResponseJSON(tcpRouter.OK, nil)
	})
	server := tcpServer.New(8080, threadPool, router)
	server.Start()
}
