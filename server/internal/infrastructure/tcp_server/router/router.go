package tcpRouter

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
)

var ErrRouteNotFound = errors.New("route not found")

type HandlerFunc func(request *RequestContext) error

type Router struct {
	routes map[RequestMeta]HandlerFunc
}

func New() *Router {
	routes := make(map[RequestMeta]HandlerFunc)
	return &Router{routes}
}

func (router *Router) AddRoute(path RequestPath, method RequestMethod, handlerFunc HandlerFunc) {
	if err := path.Validate(); err != nil {
		log.Fatal(err)
	}

	if err := method.Validate(); err != nil {
		log.Fatal(err)
	}

	rm := RequestMeta{path, method}
	router.routes[rm] = handlerFunc
	fmt.Println("New route added:", rm)
}

func (router *Router) GetHandler(meta RequestMeta) (HandlerFunc, error) {
	handler, ok := router.routes[meta]
	if !ok {
		return nil, ErrRouteNotFound
	}

	return handler, nil
}

func (router *Router) Handle(request *RequestContext) error {
	handler, err := router.GetHandler(request.RequestMeta)
	if err != nil {
		return err
	}

	return handler(request)
}

func (router *Router) ParseRequest(raw []byte) (request *RequestContext, err error) {
	err = json.Unmarshal(raw, request)
	if err != nil {
		return nil, err
	}

	return request, nil
}
