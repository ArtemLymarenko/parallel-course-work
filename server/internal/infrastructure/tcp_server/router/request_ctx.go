package tcpRouter

import (
	"encoding/json"
	"errors"
	"github.com/ArtemLymarenko/parallel-course-work/pkg/streamer"
	"net"
	"reflect"
)

type RequestContext struct {
	Request *Request
	Conn    net.Conn
}

func NewRequestContext(request *Request, conn net.Conn) *RequestContext {
	return &RequestContext{
		Request: request,
		Conn:    conn,
	}
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

func (requestCtx *RequestContext) ResponseJSON(status ResponseStatus, data any) error {
	response := Response{
		Status: status,
		Body:   data,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		return err
	}

	err = streamer.WriteBuff(requestCtx.Conn, 2048, jsonResponse)
	return err
}
