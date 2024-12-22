package tcpRouter

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
)

var ErrRouteNotFound = errors.New("route not found")

type HandlerFunc func(ctx *RequestContext) error

type Logger interface {
	Log(...interface{})
}

type Router struct {
	routes map[RequestMeta]HandlerFunc
	logger Logger
}

func New(logger Logger) *Router {
	routes := make(map[RequestMeta]HandlerFunc)
	return &Router{
		routes: routes,
		logger: logger,
	}
}

func (router *Router) AddRoute(method RequestMethod, path RequestPath, handlerFunc HandlerFunc) {
	if err := path.Validate(); err != nil {
		log.Fatal(err)
	}

	if err := method.Validate(); err != nil {
		log.Fatal(err)
	}

	rm := RequestMeta{path, method}
	router.routes[rm] = handlerFunc

	msg := fmt.Sprintf("Registered route - Method: %v, Path: %v", rm.Method, rm.Path)
	router.logger.Log(msg)
}

func (router *Router) Handle(request *Request, conn net.Conn) error {
	requestCtx := NewRequestContext(request, conn)
	handler, err := router.getHandler(requestCtx.Request.RequestMeta)
	if err != nil {
		_ = requestCtx.ResponseJSON(StatusInternalServerError, err.Error())
		return err
	}

	err = handler(requestCtx)
	if err != nil {
		_ = requestCtx.ResponseJSON(StatusInternalServerError, err.Error())
		return err
	}
	return nil
}

func (router *Router) getHandler(meta RequestMeta) (HandlerFunc, error) {
	handler, ok := router.routes[meta]
	if !ok {
		return nil, ErrRouteNotFound
	}

	return handler, nil
}

func (router *Router) ParseRawRequest(raw []byte) (*Request, error) {
	request := &Request{}
	err := json.Unmarshal(raw, request)

	if err != nil {
		return nil, err
	}

	return request, nil
}
