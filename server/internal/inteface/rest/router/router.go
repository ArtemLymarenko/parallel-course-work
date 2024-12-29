package router

import tcpRouter "server/internal/infrastructure/tcp_server/router"

type InvertedIndexHandlers interface {
	Search(ctx *tcpRouter.RequestContext) error
	AddFile(ctx *tcpRouter.RequestContext) error
	GetFileContent(ctx *tcpRouter.RequestContext) error
	RemoveFile(ctx *tcpRouter.RequestContext) error
}

func HealthCheck(ctx *tcpRouter.RequestContext) error {
	return ctx.ResponseJSON(tcpRouter.StatusOK, nil)
}

type Logger interface {
	Log(...interface{})
}

func MustInitRouter(invIndexHandlers InvertedIndexHandlers, logger Logger) *tcpRouter.Router {
	router := tcpRouter.New(logger)
	router.AddRoute(tcpRouter.GET, "/health", HealthCheck)
	router.AddRoute(tcpRouter.GET, "/index/search", invIndexHandlers.Search)
	router.AddRoute(tcpRouter.GET, "/index/file", invIndexHandlers.GetFileContent)
	router.AddRoute(tcpRouter.POST, "/index/file", invIndexHandlers.AddFile)
	router.AddRoute(tcpRouter.DELETE, "/index/file", invIndexHandlers.RemoveFile)
	return router
}
