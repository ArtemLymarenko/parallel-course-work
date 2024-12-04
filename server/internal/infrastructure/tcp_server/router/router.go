package tcpRouter

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
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
	fmt.Printf("Registered route - Method: %v, Path: %v\n", rm.Method, rm.Path)
}

func (router *Router) GetHandler(meta RequestMeta) (HandlerFunc, error) {
	handler, ok := router.routes[meta]
	fmt.Println("GetHandler", ok, handler)
	if !ok {
		return nil, ErrRouteNotFound
	}

	return handler, nil
}

func (router *Router) Handle(ctx *RequestContext) error {
	handler, err := router.GetHandler(ctx.Request.RequestMeta)
	if err != nil {
		return err
	}

	return handler(ctx)
}

func (router *Router) ParseRequest(raw []byte) (*Request, error) {
	request := &Request{}
	err := json.Unmarshal(raw, request)
	if err != nil {
		return nil, err
	}

	return request, nil
}
