package router

import tcpRouter "parallel-course-work/server/internal/infrastructure/tcp_server/router"

type InvertedIndexHandlers interface {
	Search(ctx *tcpRouter.RequestContext) error
}

func MustInitRouter(invIndexHandlers InvertedIndexHandlers) *tcpRouter.Router {
	router := tcpRouter.New()
	router.AddRoute(tcpRouter.GET, "/search", invIndexHandlers.Search)
	return router
}
