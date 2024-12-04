package tcpRouter

import (
	"encoding/json"
	"net"
)

type RequestMeta struct {
	Path   RequestPath   `json:"path"`
	Method RequestMethod `json:"method"`
}

type Response struct {
	Status ResponseStatus `json:"status"`
	Body   any            `json:"body"`
}

type Request struct {
	RequestMeta RequestMeta `json:"meta"`
	Body        any         `json:"body"`
}

type RequestContext struct {
	Conn    net.Conn
	Request *Request
}

func NewRequestContext(request *Request, conn net.Conn) *RequestContext {
	return &RequestContext{
		conn,
		request,
	}
}

func (request *RequestContext) ResponseJSON(status ResponseStatus, data any) error {
	response := Response{
		Status: status,
		Body:   data,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		return err
	}

	_, err = request.Conn.Write(jsonResponse)
	if err != nil {
		return err
	}

	return nil
}
