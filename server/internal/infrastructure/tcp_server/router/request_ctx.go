package tcpRouter

import (
	"encoding/json"
	"errors"
	"net"
	"reflect"
)

type RequestMeta struct {
	Path   RequestPath   `json:"path"`
	Method RequestMethod `json:"method"`
}

type Request struct {
	RequestMeta RequestMeta     `json:"meta"`
	Body        json.RawMessage `json:"body,omitempty"`
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

type Response struct {
	Status ResponseStatus `json:"status"`
	Body   any            `json:"body"`
}

func (requestCtx *RequestContext) ResponseJSON(status ResponseStatus, data any) error {
	response := Response{
		Status: status,
		Body:   data,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		return err
	}

	_, err = requestCtx.Conn.Write(jsonResponse)
	if err != nil {
		return err
	}

	return nil
}

func (requestCtx *RequestContext) ShouldParseBodyJSON(body any) error {
	err := json.Unmarshal(requestCtx.Request.Body, body)

	if isStructEmpty(body) {
		return errors.New("body is empty")
	}

	return err
}

func isStructEmpty(s any) bool {
	val := reflect.ValueOf(s).Elem()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)

		if !field.IsZero() {
			return false
		}
	}

	return true
}
