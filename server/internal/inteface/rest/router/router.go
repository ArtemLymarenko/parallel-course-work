package router

import tcpRouter "parallel-course-work/server/internal/infrastructure/tcp_server/router"

type InvertedIndexHandlers interface {
	Search(ctx *tcpRouter.RequestContext) error
	AddFile(ctx *tcpRouter.RequestContext) error
	GetFileContent(ctx *tcpRouter.RequestContext) error
	RemoveFile(ctx *tcpRouter.RequestContext) error
}

func HealthCheck(ctx *tcpRouter.RequestContext) error {
	return ctx.ResponseJSON(tcpRouter.StatusOK, nil)
}

func MustInitRouter(invIndexHandlers InvertedIndexHandlers) *tcpRouter.Router {
	router := tcpRouter.New()
	router.AddRoute(tcpRouter.GET, "/health", HealthCheck)
	router.AddRoute(tcpRouter.GET, "/index/search", invIndexHandlers.Search)
	router.AddRoute(tcpRouter.GET, "/index/file", invIndexHandlers.GetFileContent)
	router.AddRoute(tcpRouter.POST, "/index/file", invIndexHandlers.AddFile)
	router.AddRoute(tcpRouter.DELETE, "/index/file", invIndexHandlers.RemoveFile)
	return router
}
