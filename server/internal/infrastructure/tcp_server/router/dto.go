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

type RequestContext struct {
	conn        net.Conn
	RequestMeta RequestMeta `json:"meta"`
	Body        any         `json:"body"`
}

func (request *RequestContext) BindConn(conn net.Conn) {
	request.conn = conn
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

	_, err = request.conn.Write(jsonResponse)
	if err != nil {
		return err
	}

	return nil
}
