package tcpRouter

import (
	"encoding/json"
	"errors"
	"log"
	"net"
)

var ErrRouteNotFound = errors.New("route not found")

type HandlerFunc func(ctx *RequestContext) error

type Router struct {
	routes map[RequestMeta]HandlerFunc
}

func New() *Router {
	routes := make(map[RequestMeta]HandlerFunc)
	return &Router{routes}
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
	log.Printf("Registered route - Method: %v, Path: %v\n", rm.Method, rm.Path)
}

func (router *Router) Handle(rawRequest []byte, conn net.Conn) error {
	request, err := router.parseRawRequest(rawRequest)
	if err != nil {
		return err
	}

	requestCtx := NewRequestContext(request, conn)
	handler, err := router.getHandler(requestCtx.Request.RequestMeta)
	if err != nil {
		_ = requestCtx.ResponseJSON(InternalServerError, err.Error())
		return err
	}

	err = handler(requestCtx)
	if err != nil {
		_ = requestCtx.ResponseJSON(InternalServerError, err.Error())
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

func (router *Router) parseRawRequest(raw []byte) (*Request, error) {
	request := &Request{}
	err := json.Unmarshal(raw, request)

	if err != nil {
		return nil, err
	}

	return request, nil
}
